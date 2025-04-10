using System;
using System.Collections;
using System.Collections.Generic;
using System.Security.Cryptography;
using TMPro;
using UnityEngine;
using UnityEngine.InputSystem;
using UnityEngine.Networking;

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

    private int privateKey = 945;

    private List<PacketInfo> recievedPacketDataList;
    private bool gotPacketDataList = false;
    private List<TelemetryInfo> recievedTelemetryDataList;
    private bool gotTelemetryDataList = false;
    private List<TireInfo> recievedTireDataList;
    private bool gotTireDataList = false;

    public Dictionary<long, Packet> RecievedPackets = new Dictionary<long, Packet>();
    public Dictionary<(int, int), LapInfo> RecievedSessionLaps = new Dictionary<(int, int), LapInfo>();

    private int queryButtonCounter;

    [Header("Options")]
    public bool LogFetchOutput;

    public readonly struct PacketInfo {
        public PacketInfo(long packetID, int sessionID, int lapID, DateTime dateTime) {
            PacketID = packetID;
            SessionID = sessionID;
            LapID = lapID;
            DateTime = dateTime;
        }

        public long PacketID { get; }
        public int SessionID { get; }
        public int LapID { get; }
        public DateTime DateTime { get; }
        
        public override string ToString() => $"PacketID: {PacketID}, SessionID: {SessionID}, LapID: {LapID}, DateTime: {DateTime.ToString()}";
    }

    public readonly struct LapInfo {
        public LapInfo(int sessionID, int lapID, long lapTime, string driverName, string trackName, string trackConfiguration, string carName) {
            SessionID = sessionID;
            LapID = lapID;
            LapTime = lapTime;
            DriverName = driverName;
            TrackName = trackName;
            TrackConfiguration = trackConfiguration;
            CarName = carName;
        }

        public int SessionID { get; }
        public int LapID { get; }
        public long LapTime { get; }
        public string DriverName { get; }
        public string TrackName { get; }
        public string TrackConfiguration { get; }
        public string CarName { get; }

        public override string ToString() => $"SessionID: {SessionID}, LapID: {LapID}, LapTime: {LapTime}, DriverName: {DriverName}, TrackName: {TrackName}, TrackConfiguration: {TrackConfiguration}, CarName: {CarName}";
    }

    public readonly struct TelemetryInfo {
        public TelemetryInfo(long packetID, float speedMPH, float gas, float brake, float steer, float clutch, int gear, float rpm, float turboBoost, float localAngularVelocityX, float localAngularVelocityY, float localAngularVelocityZ, float velocityX, float velocityY, float velocityZ, float worldPositionX, float worldPositionY, float worldPositionZ, float aero_DragCoefficent, float aero_LiftCoefficentFront, float aero_LiftCoefficentRear) {
            PacketID = packetID;
            SpeedMPH = speedMPH;
            Gas = gas;
            Brake = brake;
            Steer = steer;
            Clutch = clutch;
            Gear = gear;
            RPM = rpm;
            TurboBoost = turboBoost;
            LocalAngularVelocityX = localAngularVelocityX;
            LocalAngularVelocityY = localAngularVelocityY;
            LocalAngularVelocityZ = localAngularVelocityZ;
            VelocityX = velocityX;
            VelocityY = velocityY;
            VelocityZ = velocityZ;
            WorldPositionX = worldPositionX;
            WorldPositionY = worldPositionY;
            WorldPositionZ = worldPositionZ;
            Aero_DragCoeffcient = aero_DragCoefficent;
            Aero_LiftCoefficientFront = aero_LiftCoefficentFront;
            Aero_LiftCoefficientRear = aero_LiftCoefficentRear;
        }

        public long PacketID { get; }
        public float SpeedMPH { get; }
        public float Gas { get; }
        public float Brake { get; }
        public float Steer { get; }
        public float Clutch { get; }
        public int Gear { get; }
        public float RPM { get; }
        public float TurboBoost { get; }
        public float LocalAngularVelocityX { get; }
        public float LocalAngularVelocityY { get; }
        public float LocalAngularVelocityZ { get; }
        public float VelocityX { get; }
        public float VelocityY { get; }
        public float VelocityZ { get; }
        public float WorldPositionX { get; }
        public float WorldPositionY { get; }
        public float WorldPositionZ { get; }
        public float Aero_DragCoeffcient { get; }
        public float Aero_LiftCoefficientFront { get; }
        public float Aero_LiftCoefficientRear { get; }

        public override string ToString() => $"PacketID: {PacketID}, SpeedMPH: {SpeedMPH}, Gas: {Gas}, Brake: {Brake}, Steer: {Steer}, Clutch: {Clutch}, Gear: {Gear}, RPM: {RPM}, TurboBoost: {TurboBoost}, LocalAngularVelocityX: {LocalAngularVelocityX}, LocalAngularVelocityY: {LocalAngularVelocityY}, LocalAngularVelocityZ: {LocalAngularVelocityZ}, VelocityX: {VelocityX}, VelocityY: {VelocityY}, VelocityZ: {VelocityZ}, WorldPositionX: {WorldPositionX}, WorldPositionY: {WorldPositionY}, WorldPositionZ: {WorldPositionZ}, Aero_DragCoeffcient: {Aero_DragCoeffcient}, Aero_LiftCoefficientFront: {Aero_LiftCoefficientFront}, Aero_LiftCoefficientRear: {Aero_LiftCoefficientRear}";
    }

    public readonly struct TireInfo {
        public TireInfo(long packetID, float fl_CamberRad, float fr_CamberRad, float rl_CamberRad, float rr_CamberRad, float fl_SlipAngle, float fr_SlipAngle, float rl_SlipAngle, float rr_SlipAngle, float fl_SlipRatio, float fr_SlipRatio, float rl_SlipRatio, float rr_SlipRatio, float fl_SelfAligningTorque, float fr_SelfAligningTorque, float rl_SelfAligningTorque, float rr_SelfAligningTorque, float fl_Load, float fr_Load, float rl_Load, float rr_Load, float fl_TyreSlip, float fr_TyreSlip, float rl_TyreSlip, float rr_TyreSlip, float fl_ThermalState, float fr_ThermalState, float rl_ThermalState, float rr_ThermalState, float fl_DynamicPressure, float fr_DynamicPressure, float rl_DynamicPressure, float rr_DynamicPressure, float fl_TyreDirtyLevel, float fr_TyreDirtyLevel, float rl_TyreDirtyLevel, float rr_TyreDirtyLevel) {
            PacketID = packetID;
            FL_CamberRad = fl_CamberRad;
            FR_CamberRad = fr_CamberRad;
            RL_CamberRad = rl_CamberRad;
            RR_CamberRad = rr_CamberRad;
            FL_SlipAngle = fl_SlipAngle;
            FR_SlipAngle = fr_SlipAngle;
            RL_SlipAngle = rl_SlipAngle;
            RR_SlipAngle = rr_SlipAngle;
            FL_SlipRatio = fl_SlipRatio;
            FR_SlipRatio = fr_SlipRatio;
            RL_SlipRatio = rl_SlipRatio;
            RR_SlipRatio = rr_SlipRatio;
            FL_SelfAligningTorque = fl_SelfAligningTorque;
            FR_SelfAligningTorque = fr_SelfAligningTorque;
            RL_SelfAligningTorque = rl_SelfAligningTorque;
            RR_SelfAligningTorque = rr_SelfAligningTorque;
            FL_Load = fl_Load;
            FR_Load = fr_Load;
            RL_Load = rl_Load;
            RR_Load = rr_Load;
            FL_TyreSlip = fl_TyreSlip;
            FR_TyreSlip= fr_TyreSlip;
            RL_TyreSlip = rl_TyreSlip;
            RR_TyreSlip = rr_TyreSlip;
            FL_ThermalState = fl_ThermalState;
            FR_ThermalState = fr_ThermalState;
            RL_ThermalState = rl_ThermalState;
            RR_ThermalState = rl_ThermalState;
            FL_DynamicPressure = fl_DynamicPressure;
            FR_DynamicPressure = fr_DynamicPressure;
            RL_DynamicPressure = rl_DynamicPressure;
            RR_DynamicPressure = rr_DynamicPressure;
            FL_TyreDirtyLevel = fl_TyreDirtyLevel;
            FR_TyreDirtyLevel = fr_TyreDirtyLevel;
            RL_TyreDirtyLevel = rl_TyreDirtyLevel;
            RR_TyreDirtyLevel = rr_TyreDirtyLevel;
        }

        public long PacketID { get; }
        public float FL_CamberRad { get; }
        public float FR_CamberRad { get; }
        public float RL_CamberRad { get; }
        public float RR_CamberRad { get; }
        public float FL_SlipAngle { get; }
        public float FR_SlipAngle { get; }
        public float RL_SlipAngle { get; }
        public float RR_SlipAngle { get; }
        public float FL_SlipRatio { get; }
        public float FR_SlipRatio { get; }
        public float RL_SlipRatio { get; }
        public float RR_SlipRatio { get; }
        public float FL_SelfAligningTorque { get; }
        public float FR_SelfAligningTorque { get; }
        public float RL_SelfAligningTorque { get; }
        public float RR_SelfAligningTorque { get; }
        public float FL_Load { get; }
        public float FR_Load { get; }
        public float RL_Load { get; }
        public float RR_Load { get; }
        public float FL_TyreSlip { get; }
        public float FR_TyreSlip { get; }
        public float RL_TyreSlip { get; }
        public float RR_TyreSlip { get; }
        public float FL_ThermalState { get; }
        public float FR_ThermalState { get; }
        public float RL_ThermalState { get; }
        public float RR_ThermalState { get; }
        public float FL_DynamicPressure { get; }
        public float FR_DynamicPressure { get; }
        public float RL_DynamicPressure { get; }
        public float RR_DynamicPressure { get; }
        public float FL_TyreDirtyLevel { get; }
        public float FR_TyreDirtyLevel { get; }
        public float RL_TyreDirtyLevel { get; }
        public float RR_TyreDirtyLevel { get; }

        public override string ToString() => $"PacketID: {PacketID}, FL_CamberRad: {FL_CamberRad}, FR_CamberRad: {FR_CamberRad}, RL_CamberRad: {RL_CamberRad}, RR_CamberRad: {RR_CamberRad}, FL_SlipAngle: {FL_SlipAngle}, FR_SlipAngle: {FR_SlipAngle}, RL_SlipAngle: {RL_SlipAngle}, RR_SlipAngle: {RR_SlipAngle}, FL_SlipRatio: {FL_SlipRatio}, FR_SlipRatio: {FR_SlipRatio}, RL_SlipRatio: {RL_SlipRatio}, RR_SlipRatio: {RR_SlipRatio}, FL_SelfAligningTorque: {FL_SelfAligningTorque}, FR_SelfAligningTorque: {FR_SelfAligningTorque}, RL_SelfAligningTorque: {RL_SelfAligningTorque}, RR_SelfAligningTorque: {RR_SelfAligningTorque}, FL_Load: {FL_Load}, FR_Load: {FR_Load}, RL_Load: {RL_Load}, RR_Load: {RR_Load}, FL_TyreSlip: {FL_TyreSlip}, FR_TyreSlip: {FR_TyreSlip}, RL_TyreSlip: {RL_TyreSlip}, RR_TyreSlip: {RR_TyreSlip}, FL_ThermalState: {FL_ThermalState}, FR_ThermalState: {FR_ThermalState}, RL_ThermalState: {RL_ThermalState}, RR_ThermalState: {RR_ThermalState}, FL_DynamicPressure: {FL_DynamicPressure}, FR_DynamicPressure: {FR_DynamicPressure}, RL_DynamicPressure: {RL_DynamicPressure}, RR_DynamicPressure: {RR_DynamicPressure}, FL_TyreDirtyLevel: {FL_TyreDirtyLevel}, FR_TyreDirtyLevel: {FR_TyreDirtyLevel}, RL_TyreDirtyLevel: {RL_TyreDirtyLevel}, RR_TyreDirtyLevel: {RR_TyreDirtyLevel}";
    }

    public readonly struct Packet {
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
        StartCoroutine(getRequest(getQueryString(1, 0, 100)));
        StartCoroutine(getRequest(getQueryString(3, 0, 100)));
        StartCoroutine(getRequest(getQueryString(4, 0, 100)));
    }

    public void InsertDataButtonClick() {
        StartCoroutine(randomInsert());
    }

    public void QueryPackets(long packetIDStart, long packetIDEnd) {
        StartCoroutine(getRequest(getQueryString(1, packetIDStart, packetIDEnd)));
        StartCoroutine(getRequest(getQueryString(3, packetIDStart, packetIDEnd)));
        StartCoroutine(getRequest(getQueryString(4, packetIDStart, packetIDEnd)));
    }

    public void QuerySessionLap(int sessionID, int lapID) {
        StartCoroutine(getRequest(getQueryString(2, sessionID, lapID)));
    }

    void Awake() {
        if (Instance != null && Instance != this) {
            Destroy(this);
        } else {
            Instance = this;
        }
    }

    void Start() {
        terminalStart();
    }

    void Update() {
        terminalUpdate();

        if (gotPacketDataList && gotTelemetryDataList && gotTireDataList) {
            gotPacketDataList = false;
            gotTelemetryDataList = false;
            gotTireDataList = false;

            if (recievedPacketDataList.Count == recievedTelemetryDataList.Count && recievedTelemetryDataList.Count == recievedTireDataList.Count) {
                for (int i = 0; i < recievedPacketDataList.Count; i++) {
                    PacketInfo packetData = recievedPacketDataList[i];
                    TelemetryInfo telemetryData = recievedTelemetryDataList[i];
                    TireInfo tireData = recievedTireDataList[i];
                    RecievedPackets[packetData.PacketID] = new Packet(packetData, telemetryData, tireData);
                    if (terminalActive && packetData.PacketID == queryButtonCounter) {
                        log($"{RecievedPackets[queryButtonCounter++]}\n");
                    }
                }
            } else {
                Debug.Log($"Error Got Mismatched Lengths packetDataList {recievedPacketDataList.Count} items, telemetryDataList {recievedTelemetryDataList.Count} items, tireDataList {recievedTireDataList.Count} items");
            }
        }
    }

    private IEnumerator getRequest(string uri) {
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
                    string[] packetStrings = webRequest.downloadHandler.text.Split("\n");

                    if (LogFetchOutput) { 
                        Debug.Log($"GET: {string.Join(", ", packetStrings)}"); 
                    }

                    if (packetStrings.Length > 0) {
                        switch (packetStrings[0].Split(",").Length) {
                            case 4:
                                List<PacketInfo> packetDataList = new List<PacketInfo>(); 
                                foreach (string packet in packetStrings) {
                                    string[] packetArray = packet.Split(",");
                                    if (packetArray.Length != 4) { continue; } 
                                    PacketInfo packetData = new PacketInfo(long.Parse(packetArray[0]), int.Parse(packetArray[1]), int.Parse(packetArray[2]), DateTime.Parse(packetArray[3])); 
                                    packetDataList.Add(packetData);
                                }
                                recievedPacketDataList = packetDataList;
                                gotPacketDataList = true;
                                break;
                            case 21:
                                List<TelemetryInfo> telemetryDataList = new List<TelemetryInfo>(); 
                                foreach (string packet in packetStrings) {
                                    string[] packetArray = packet.Split(",");
                                    if (packetArray.Length != 21) { continue; }
                                    TelemetryInfo telemetryData = new TelemetryInfo(long.Parse(packetArray[0]), float.Parse(packetArray[1]), float.Parse(packetArray[2]), float.Parse(packetArray[3]), float.Parse(packetArray[4]), float.Parse(packetArray[5]), int.Parse(packetArray[6]), float.Parse(packetArray[7]), float.Parse(packetArray[8]), float.Parse(packetArray[9]), float.Parse(packetArray[10]), float.Parse(packetArray[11]), float.Parse(packetArray[12]), float.Parse(packetArray[13]), float.Parse(packetArray[14]), float.Parse(packetArray[15]), float.Parse(packetArray[16]), float.Parse(packetArray[17]), float.Parse(packetArray[18]), float.Parse(packetArray[19]), float.Parse(packetArray[20])); 
                                    telemetryDataList.Add(telemetryData);
                                }
                                recievedTelemetryDataList = telemetryDataList;
                                gotTelemetryDataList = true;
                                break;
                            case 37:
                                List<TireInfo> tireDataList = new List<TireInfo>(); 
                                foreach (string packet in packetStrings) {
                                    string[] packetArray = packet.Split(",");
                                    if (packetArray.Length != 37) { continue; }
                                    TireInfo tireData = new TireInfo(long.Parse(packetArray[0]), float.Parse(packetArray[1]), float.Parse(packetArray[2]), float.Parse(packetArray[3]), float.Parse(packetArray[4]), float.Parse(packetArray[5]), float.Parse(packetArray[6]), float.Parse(packetArray[7]), float.Parse(packetArray[8]), float.Parse(packetArray[9]), float.Parse(packetArray[10]), float.Parse(packetArray[11]), float.Parse(packetArray[12]), float.Parse(packetArray[13]), float.Parse(packetArray[14]), float.Parse(packetArray[15]), float.Parse(packetArray[16]), float.Parse(packetArray[17]), float.Parse(packetArray[18]), float.Parse(packetArray[19]), float.Parse(packetArray[20]), float.Parse(packetArray[21]), float.Parse(packetArray[22]), float.Parse(packetArray[23]), float.Parse(packetArray[24]), float.Parse(packetArray[25]), float.Parse(packetArray[26]), float.Parse(packetArray[27]), float.Parse(packetArray[28]), float.Parse(packetArray[29]), float.Parse(packetArray[30]), float.Parse(packetArray[31]), float.Parse(packetArray[32]), float.Parse(packetArray[33]), float.Parse(packetArray[34]), float.Parse(packetArray[35]), float.Parse(packetArray[36])); 
                                    tireDataList.Add(tireData);
                                }
                                recievedTireDataList = tireDataList;
                                gotTireDataList = true;
                                break;
                            case 7:
                                foreach (string packet in packetStrings) {
                                    string[] packetArray = packet.Split(",");
                                    if (packetArray.Length != 7) { continue; }
                                    LapInfo lapData = new LapInfo(int.Parse(packetArray[0]), int.Parse(packetArray[1]), long.Parse(packetArray[2]), packetArray[3].Replace("{Comma}", ","), packetArray[4].Replace("{Comma}", ","), packetArray[5].Replace("{Comma}", ","), packetArray[6].Replace("{Comma}", ","));
                                    RecievedSessionLaps[(lapData.SessionID, lapData.LapID)] = lapData;
                                }
                                break;
                        }
                    }
                    break;
            }
        }
    }

    private IEnumerator randomInsert() {
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

    private string getQueryString(int queryType, long start, long end) {
        //Query Types
        //1 = query PacketInfo from PackedID to PacketID
        //2 = query LapInfo from LapID to LapID
        //3 = query TelemetryInfo from PacketID to PacketID
        //4 = query TireInfo from PacketID to PacketID

        switch (queryType) {
            case 1:
                return $"http://localhost:8080/sqliteQuery?table=PacketInfo&start={start}&end={end}";
            case 2:
                return $"http://localhost:8080/sqliteQuery?table=LapInfo&start={start}&end={end}";
            case 3:
                return $"http://localhost:8080/sqliteQuery?table=TelemetryInfo&start={start}&end={end}";
            case 4:
                return $"http://localhost:8080/sqliteQuery?table=TireInfo&start={start}&end={end}";
            default:
                throw new Exception($"queryType of {queryType} unknown");
        }
    }

    private string hashInput(string input) {
        SHA256Managed hm = new SHA256Managed();
        byte[] hashValue = hm.ComputeHash(System.Text.Encoding.ASCII.GetBytes(input));
        string hash_convert = BitConverter.ToString(hashValue).Replace("-", "").ToLower();
        return hash_convert;
    }

    private void terminalStart() {
        outputTextBox = Terminal.transform.GetChild(0).GetChild(0).GetComponent<TextMeshProUGUI>();

        terminalLineP = 0;
        terminalLines = new List<string>(){""};

        scrollAction = new InputAction("Scroll", binding: "<Mouse>/scroll");
        scrollAction.Enable();
        scrollAction.performed += x => { terminalLineP += (int)x.ReadValue<Vector2>()[1]; terminalStateChanged = true; };
        terminalActive = true;

        queryButtonCounter = 1;
    }

    private void terminalUpdate() {
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

    private void clearLog() {
        if (terminalActive) {
            terminalLines = new List<string>(){""};
            terminalLineP = 0;

            terminalStateChanged = true;
        }
    }

    private void log(string text) {
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