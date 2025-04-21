using System.Collections.Generic;
using CarSpeedometerPro.Speedometer;
using UnityEngine;

public class MainManager : MonoBehaviour
{
    public static MainManager Instance { get; private set; }
    private DatabaseAccess db;

    [HideInInspector]
    public long SelectedPacketID = 0;
    [HideInInspector]
    public DatabaseAccess.TelemetryInfo SelectedTelemetryInfo;

    [Header("Packet Range")]
    public long StartPacketID = 0;
    public long EndPacketID = 2546;
    

    private float localTime = 0;

    [Space]
    [Header("Playback")]
    public float PlayingStepRate = 0.1f;
    private float playingStep = 0;
    public bool Paused = false;

    [Space]
    [Header("UI Refs")]
    public SpeedometerController speedometer;
    public TrackMap trackMap;
    public CarPosition car;
    public LineRenderer SteerGraph;
    public LineRenderer ThrottleGraph;
    public LineRenderer BrakeGraph;
    public LineRenderer SpeedGraph;

    private bool setTrackMap = false;

    void Awake() {
        if (Instance != null && Instance != this) {
            Destroy(this);
        } else {
            Instance = this;
        }

    }

    void Start() {
        db = DatabaseAccess.Instance;
        db.QueryTelemetryInfo(StartPacketID, EndPacketID);
    }

    void Update() {
        if (!setTrackMap && db.TelemetryInfoQueried) {
            setTrackMap = true;
            float[] trackData = new float[(EndPacketID-StartPacketID +1)*2];
            int counter = 0;
            for (long i = StartPacketID; i <= EndPacketID; i++) {
                trackData[counter++] = db.RecievedTelemetryInfo[i].WorldPositionX;
                trackData[counter++] = db.RecievedTelemetryInfo[i].WorldPositionZ;
            }
            trackMap.TrackData = trackData;
        }

        localTime += Time.deltaTime;
        if (!Paused && localTime > playingStep && db.TelemetryInfoQueried) {
            playingStep = localTime + PlayingStepRate;
            SelectedPacketID += 1;

            if (SelectedPacketID > EndPacketID) { SelectedPacketID = StartPacketID; }
            
            SelectedTelemetryInfo = DatabaseAccess.Instance.RecievedTelemetryInfo.ContainsKey(SelectedPacketID) ? DatabaseAccess.Instance.RecievedTelemetryInfo[SelectedPacketID] : null;
            
            if (SelectedTelemetryInfo == null) { return; }

            speedometer.throttleSliderValue = SelectedTelemetryInfo.Gas;
            speedometer.brakeSliderValue = SelectedTelemetryInfo.Brake;
            speedometer.currentSpeedText.text = $"{Mathf.RoundToInt(SelectedTelemetryInfo.SpeedMPH)}";

            float x = (SelectedTelemetryInfo.WorldPositionX - trackMap.minX) * trackMap.scaleFactor + trackMap.centerOffsetX + trackMap.ManualOffset.x;
            float y = (SelectedTelemetryInfo.WorldPositionZ - trackMap.minY) * trackMap.scaleFactor + trackMap.centerOffsetY + trackMap.ManualOffset.y;

            car.MovePosition(new Vector2(x, y));

            //HEY BTW THE TIME RN FROM THE DATA IS FUCKED SINCE I DIDNT HAVE TIME (HAHAH GET IT) TO GET IT WORKING. DO NOT USE THE DATE/TIME THINGY AS IT WILL GIVE INACCURATE NUMBERS

            int graphNodes = 50;

            List<Vector2> steerPoints = new List<Vector2>();
            float counter = -4;
            if ((SelectedPacketID - StartPacketID) < graphNodes/2) {
                for (long i = StartPacketID; i <= SelectedPacketID; i++) {
                    steerPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Steer));
                }
            } else if ((EndPacketID - SelectedPacketID) < graphNodes/2) {
                for (long i = SelectedPacketID-(graphNodes/2); i <= EndPacketID; i++) {
                    steerPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Steer));
                }
            } else {
                for (long i = SelectedPacketID-(graphNodes/2); i <= SelectedPacketID+(graphNodes/2); i++) {
                    steerPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Steer));
                }
            }
            SteerGraph.SetPoints(steerPoints);

            List<Vector2> throttlePoints = new List<Vector2>();
            counter = -4;
            if ((SelectedPacketID - StartPacketID) < graphNodes/2) {
                for (long i = StartPacketID; i <= SelectedPacketID; i++) {
                    throttlePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Gas));
                }
            } else if ((EndPacketID - SelectedPacketID) < graphNodes/2) {
                for (long i = SelectedPacketID-(graphNodes/2); i <= EndPacketID; i++) {
                    throttlePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Gas));
                }
            } else {
                for (long i = SelectedPacketID-(graphNodes/2); i <= SelectedPacketID+(graphNodes/2); i++) {
                    throttlePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Gas));
                }
            }
            ThrottleGraph.SetPoints(throttlePoints);

            List<Vector2> brakePoints = new List<Vector2>();
            counter = -4;
            if ((SelectedPacketID - StartPacketID) < graphNodes/2) {
                for (long i = StartPacketID; i <= SelectedPacketID; i++) {
                    brakePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Brake));
                }
            } else if ((EndPacketID - SelectedPacketID) < graphNodes/2) {
                for (long i = SelectedPacketID-(graphNodes/2); i <= EndPacketID; i++) {
                    brakePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Brake));
                }
            } else {
                for (long i = SelectedPacketID-(graphNodes/2); i <= SelectedPacketID+(graphNodes/2); i++) {
                    brakePoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].Brake));
                }
            }
            BrakeGraph.SetPoints(brakePoints);

            List<Vector2> speedPoints = new List<Vector2>();
            counter = -4;
            if ((SelectedPacketID - StartPacketID) < graphNodes/2) {
                for (long i = StartPacketID; i <= SelectedPacketID; i++) {
                    speedPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].SpeedMPH));
                }
            } else if ((EndPacketID - SelectedPacketID) < graphNodes/2) {
                for (long i = SelectedPacketID-(graphNodes/2); i <= EndPacketID; i++) {
                    speedPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].SpeedMPH));
                }
            } else {
                for (long i = SelectedPacketID-(graphNodes/2); i <= SelectedPacketID+(graphNodes/2); i++) {
                    speedPoints.Add(new Vector2(counter++, db.RecievedTelemetryInfo[i].SpeedMPH));
                }
            }
            SpeedGraph.SetPoints(speedPoints);
        }
    }
}
