using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using TMPro;
using UnityEngine;
using UnityEngine.InputSystem;
using UnityEngine.Networking;
using Debug = UnityEngine.Debug;

public class DatabaseAccess : MonoBehaviour {
    public static DatabaseAccess Instance { get; private set; }
    
    public Canvas Terminal;
    private TextMeshProUGUI outputTextBox;

    private InputAction scrollAction;

    private int terminalLineLimit = 21;
    private bool terminalStateChanged;
    public int terminalLineP;
    private List<string> terminalLines;
    private bool terminalActive = false;

    private PacketInfo[] recievedPacketDataList;
    private bool gotPacketDataList = false;
    private TelemetryInfo[] recievedTelemetryDataList;
    private bool gotTelemetryDataList = false;
    private TireInfo[] recievedTireDataList;
    private bool gotTireDataList = false;

    public Dictionary<long, Packet> RecievedPackets = new Dictionary<long, Packet>();
    public Dictionary<(int, int), LapInfo> RecievedSessionLaps = new Dictionary<(int, int), LapInfo>();

    [Header("Options")]
    public bool LogFetchOutput;

    [Serializable]
    public class PacketInfo {
        public long PacketID;
        public int SessionID;
        public int LapID;
        public string PacketDatetime;
        
        public override string ToString() => $"PacketID: {PacketID}, SessionID: {SessionID}, LapID: {LapID}, DateTime: {PacketDatetime}";
    }

    [Serializable]
    public class LapInfo {
        public int SessionID;
        public int LapID;
        public long LapTime;
        public string DriverName;
        public string TrackName;
        public string TrackConfiguration;
        public string CarName;

        public override string ToString() => $"SessionID: {SessionID}, LapID: {LapID}, LapTime: {LapTime}, DriverName: {DriverName}, TrackName: {TrackName}, TrackConfiguration: {TrackConfiguration}, CarName: {CarName}";
    }

    [Serializable]
    public class TelemetryInfo {
        public long PacketID;
        public float SpeedMPH;
        public float Gas;
        public float Brake;
        public float Steer;
        public float Clutch;
        public int Gear;
        public float RPM;
        public float TurboBoost;
        public float LocalAngularVelocityX;
        public float LocalAngularVelocityY;
        public float LocalAngularVelocityZ;
        public float VelocityX;
        public float VelocityY;
        public float VelocityZ;
        public float WorldPositionX;
        public float WorldPositionY;
        public float WorldPositionZ;
        public float Aero_DragCoeffcient;
        public float Aero_LiftCoefficientFront;
        public float Aero_LiftCoefficientRear;

        public override string ToString() => $"PacketID: {PacketID}, SpeedMPH: {SpeedMPH}, Gas: {Gas}, Brake: {Brake}, Steer: {Steer}, Clutch: {Clutch}, Gear: {Gear}, RPM: {RPM}, TurboBoost: {TurboBoost}, LocalAngularVelocityX: {LocalAngularVelocityX}, LocalAngularVelocityY: {LocalAngularVelocityY}, LocalAngularVelocityZ: {LocalAngularVelocityZ}, VelocityX: {VelocityX}, VelocityY: {VelocityY}, VelocityZ: {VelocityZ}, WorldPositionX: {WorldPositionX}, WorldPositionY: {WorldPositionY}, WorldPositionZ: {WorldPositionZ}, Aero_DragCoeffcient: {Aero_DragCoeffcient}, Aero_LiftCoefficientFront: {Aero_LiftCoefficientFront}, Aero_LiftCoefficientRear: {Aero_LiftCoefficientRear}";
    }

    [Serializable]
    public class TireInfo {
        public long PacketID;
        public float FL_CamberRad;
        public float FR_CamberRad;
        public float RL_CamberRad;
        public float RR_CamberRad;
        public float FL_SlipAngle;
        public float FR_SlipAngle;
        public float RL_SlipAngle;
        public float RR_SlipAngle;
        public float FL_SlipRatio;
        public float FR_SlipRatio;
        public float RL_SlipRatio;
        public float RR_SlipRatio;
        public float FL_SelfAligningTorque;
        public float FR_SelfAligningTorque;
        public float RL_SelfAligningTorque;
        public float RR_SelfAligningTorque;
        public float FL_Load;
        public float FR_Load;
        public float RL_Load;
        public float RR_Load;
        public float FL_TyreSlip;
        public float FR_TyreSlip;
        public float RL_TyreSlip;
        public float RR_TyreSlip;
        public float FL_ThermalState;
        public float FR_ThermalState;
        public float RL_ThermalState;
        public float RR_ThermalState;
        public float FL_DynamicPressure;
        public float FR_DynamicPressure;
        public float RL_DynamicPressure;
        public float RR_DynamicPressure;
        public float FL_TyreDirtyLevel;
        public float FR_TyreDirtyLevel;
        public float RL_TyreDirtyLevel;
        public float RR_TyreDirtyLevel;

        public override string ToString() => $"PacketID: {PacketID}, FL_CamberRad: {FL_CamberRad}, FR_CamberRad: {FR_CamberRad}, RL_CamberRad: {RL_CamberRad}, RR_CamberRad: {RR_CamberRad}, FL_SlipAngle: {FL_SlipAngle}, FR_SlipAngle: {FR_SlipAngle}, RL_SlipAngle: {RL_SlipAngle}, RR_SlipAngle: {RR_SlipAngle}, FL_SlipRatio: {FL_SlipRatio}, FR_SlipRatio: {FR_SlipRatio}, RL_SlipRatio: {RL_SlipRatio}, RR_SlipRatio: {RR_SlipRatio}, FL_SelfAligningTorque: {FL_SelfAligningTorque}, FR_SelfAligningTorque: {FR_SelfAligningTorque}, RL_SelfAligningTorque: {RL_SelfAligningTorque}, RR_SelfAligningTorque: {RR_SelfAligningTorque}, FL_Load: {FL_Load}, FR_Load: {FR_Load}, RL_Load: {RL_Load}, RR_Load: {RR_Load}, FL_TyreSlip: {FL_TyreSlip}, FR_TyreSlip: {FR_TyreSlip}, RL_TyreSlip: {RL_TyreSlip}, RR_TyreSlip: {RR_TyreSlip}, FL_ThermalState: {FL_ThermalState}, FR_ThermalState: {FR_ThermalState}, RL_ThermalState: {RL_ThermalState}, RR_ThermalState: {RR_ThermalState}, FL_DynamicPressure: {FL_DynamicPressure}, FR_DynamicPressure: {FR_DynamicPressure}, RL_DynamicPressure: {RL_DynamicPressure}, RR_DynamicPressure: {RR_DynamicPressure}, FL_TyreDirtyLevel: {FL_TyreDirtyLevel}, FR_TyreDirtyLevel: {FR_TyreDirtyLevel}, RL_TyreDirtyLevel: {RL_TyreDirtyLevel}, RR_TyreDirtyLevel: {RR_TyreDirtyLevel}";
    }

    public class Packet {
        public Packet(PacketInfo packetData, TelemetryInfo telemetryData, TireInfo tireData) {
            PacketData = packetData;
            TelemetryData = telemetryData;
            TireData = tireData;
        }

        public PacketInfo PacketData { get; }
        public TelemetryInfo TelemetryData { get; }
        public TireInfo TireData { get; }

        public override string ToString() => $"PacketData:  {PacketData}\nTelemetryData:    {TelemetryData}\nTireData:  {TireData}";
    }

    public void GetQueryButtonClick() {
        QueryPackets(0, 100);
    }

    public void InsertDataButtonClick() {
        StartCoroutine(RandomInsert());
    }

    public void QueryPackets(long packetIDStart, long packetIDEnd) {
        StartCoroutine(QueryRequest(1, packetIDStart, packetIDEnd));
        StartCoroutine(QueryRequest(3, packetIDStart, packetIDEnd));
        StartCoroutine(QueryRequest(4, packetIDStart, packetIDEnd));
    }

    public void QuerySessionLap(int sessionID, int lapID) {
        StartCoroutine(QueryRequest(2, sessionID, lapID));
    }

    void Awake() {
        if (Instance != null && Instance != this) {
            Destroy(this);
        } else {
            Instance = this;
        }

        Process golangServer = new Process();
        golangServer.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
        golangServer.StartInfo.FileName = "goServer.exe";
        golangServer.Start();
    }

    void Start() {
        TerminalStart();
    }

    void Update() {
        TerminalUpdate();

        if (gotPacketDataList && gotTelemetryDataList && gotTireDataList) {
            gotPacketDataList = false;
            gotTelemetryDataList = false;
            gotTireDataList = false;
            
            if (recievedPacketDataList.Length == recievedTelemetryDataList.Length && recievedTelemetryDataList.Length == recievedTireDataList.Length) {
                for (int i = 0; i < recievedPacketDataList.Length; i++) {
                    PacketInfo packetData = recievedPacketDataList[i];
                    TelemetryInfo telemetryData = recievedTelemetryDataList[i];
                    TireInfo tireData = recievedTireDataList[i];
                    RecievedPackets[packetData.PacketID] = new Packet(packetData, telemetryData, tireData);
                    if (terminalActive) {
                        Log($"{RecievedPackets[packetData.PacketID]}\n");
                    }
                }
            } else {
                Debug.Log($"Error Got Mismatched Lengths packetDataList {recievedPacketDataList.Length} items, telemetryDataList {recievedTelemetryDataList.Length} items, tireDataList {recievedTireDataList.Length} items");
            }
        }
    }

    private IEnumerator QueryRequest(int queryType, long start, long end) {
        //Query Types
        //1 = query PacketInfo from PackedID to PacketID
        //2 = query LapInfo from LapID to LapID
        //3 = query TelemetryInfo from PacketID to PacketID
        //4 = query TireInfo from PacketID to PacketID

        string uri;

        switch (queryType) {
            case 1:
                uri = $"http://localhost:8080/sqliteQuery?table=PacketInfo&start={start}&end={end}";
                break;
            case 2:
                uri = $"http://localhost:8080/sqliteQuery?table=LapInfo&start={start}&end={end}";
                break;
            case 3:
                uri = $"http://localhost:8080/sqliteQuery?table=TelemetryInfo&start={start}&end={end}";
                break;
            case 4:
                uri = $"http://localhost:8080/sqliteQuery?table=TireInfo&start={start}&end={end}";
                break;
            default:
                throw new Exception($"queryType of {queryType} unknown");
        }

        using (UnityWebRequest webRequest = UnityWebRequest.Get(uri)) {
            yield return webRequest.SendWebRequest();
  
            string[] pages = uri.Split('/');
            int page = pages.Length - 1;

            switch (webRequest.result) {
                case UnityWebRequest.Result.ConnectionError:
                    Debug.LogError("Connection Error: " + webRequest.error);
                    break;
                case UnityWebRequest.Result.DataProcessingError:
                    Debug.LogError(pages[page] + ": Error: " + webRequest.error);
                    break;
                case UnityWebRequest.Result.ProtocolError:
                    Debug.LogError(pages[page] + ": HTTP Error: " + webRequest.error);
                    break;
                case UnityWebRequest.Result.Success:
                    string[] infoBits = webRequest.downloadHandler.text.Substring(1, webRequest.downloadHandler.text.Length-3).Split("},");
                    switch (queryType) {
                        case 1:
                            recievedPacketDataList = new PacketInfo[infoBits.Length];
                            for (int i = 0; i < infoBits.Length; i++) {
                                PacketInfo packetInfo;
                                if (i == infoBits.Length -1) { 
                                    packetInfo = recievedPacketDataList[i] = JsonUtility.FromJson<PacketInfo>(infoBits[i]); 
                                } else {
                                    packetInfo = JsonUtility.FromJson<PacketInfo>(infoBits[i] + "}");
                                }
                                recievedPacketDataList[i] = packetInfo;
                            }
                            gotPacketDataList = true; 
                            break;
                        case 2:
                            for (int i = 0; i < infoBits.Length; i++) {
                                LapInfo lapInfo;
                                if (i == infoBits.Length -1) { 
                                    lapInfo = JsonUtility.FromJson<LapInfo>(infoBits[i]); 
                                } else {
                                    lapInfo = JsonUtility.FromJson<LapInfo>(infoBits[i] + "}");
                                }
                                RecievedSessionLaps[(lapInfo.SessionID, lapInfo.LapID)] = lapInfo;
                            }
                            break;
                        case 3:
                            recievedTelemetryDataList = new TelemetryInfo[infoBits.Length];
                            for (int i = 0; i < infoBits.Length; i++) {
                                TelemetryInfo telemetryInfo;
                                if (i == infoBits.Length -1) { 
                                    telemetryInfo = JsonUtility.FromJson<TelemetryInfo>(infoBits[i]); 
                                } else {
                                    telemetryInfo = JsonUtility.FromJson<TelemetryInfo>(infoBits[i] + "}");
                                }
                                recievedTelemetryDataList[i] = telemetryInfo;
                            }
                            gotTelemetryDataList = true; 
                            break;
                        case 4:
                            recievedTireDataList = new TireInfo[infoBits.Length];
                            for (int i = 0; i < infoBits.Length; i++) {
                                TireInfo tireInfo;
                                if (i == infoBits.Length -1) { 
                                    tireInfo = JsonUtility.FromJson<TireInfo>(infoBits[i]); 
                                } else {
                                    tireInfo = JsonUtility.FromJson<TireInfo>(infoBits[i] + "}");
                                }
                                recievedTireDataList[i] = tireInfo;
                            }
                            gotTireDataList = true; 
                            break;
                    } 
                    break;
            }
        }
    }

    private IEnumerator RandomInsert() {
        using (UnityWebRequest request = new UnityWebRequest("http://localhost:8080/csvInsert", UnityWebRequest.kHttpVerbPOST)) {
            request.SetRequestHeader("Content-Type", "test/csv");
            request.uploadHandler = new UploadHandlerFile("./telemetryData/RandomInsertData.csv");

            yield return request.SendWebRequest();

            switch (request.result) {
            case UnityWebRequest.Result.ConnectionError:
                Debug.LogError("Connection Error: " + request.error);
                break;
            case UnityWebRequest.Result.DataProcessingError:
                Debug.LogError("Error: " + request.error);
                break;
            case UnityWebRequest.Result.ProtocolError:
                Debug.LogError("HTTP Error: " + request.error);
                break;
            case UnityWebRequest.Result.Success:
                Debug.Log("Post Success: " + request.result);
                break;
            }
        }
    }

    private void TerminalStart() {
        outputTextBox = Terminal.transform.GetChild(0).GetChild(0).GetComponent<TextMeshProUGUI>();

        terminalLineP = 0;
        terminalLines = new List<string>(){""};

        scrollAction = new InputAction("Scroll", binding: "<Mouse>/scroll");
        scrollAction.Enable();
        scrollAction.performed += x => { terminalLineP += (int)x.ReadValue<Vector2>()[1]; terminalStateChanged = true; };
        terminalActive = true;
    }

    private void TerminalUpdate() {
        if (terminalActive && terminalStateChanged) {
            terminalStateChanged = false;

            if (terminalLineP > terminalLines.Count) {
                terminalLineP = terminalLines.Count;
            } else if (terminalLineP < terminalLineLimit) {
                terminalLineP = terminalLineLimit;
            }

            string newTerminalText = "";
            for (int i = terminalLineP - terminalLineLimit; i <= terminalLineP; i++) {
                if (i >= 0 && i < terminalLines.Count) { newTerminalText += terminalLines[i] + "\n"; }
            }
            outputTextBox.text = newTerminalText;
        }
    }

    public void ClearLog() {
        if (terminalActive) {
            terminalLines = new List<string>(){""};
            terminalLineP = 0;

            terminalStateChanged = true;
        }
    }

    public void Log(string text) {
        if (terminalActive) {
            text.Replace("\r\n", "\n");
            text.Replace("<br>", "\n");
            if (text.Contains("\n")) {
                var lines = text.Split("\n");
                foreach(string line in lines) {
                    terminalLines[terminalLines.Count - 1] += line.Replace("\n", "");
                    terminalLines.Add("");
                }
            } else {
                terminalLines[terminalLines.Count - 1] += text;
            }

            terminalLineP = terminalLines.Count;

            terminalStateChanged = true;
        }   
    }
}