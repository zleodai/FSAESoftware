using System;
using UnityEngine;

public class MainManager : MonoBehaviour
{
    public static MainManager Instance { get; private set; }


    public long SelectedPacketID = 0;

    public DatabaseAccess.Packet SelectedPacket;
    
    public bool Paused = true;

    private float localTime = 0;
    
    public float PlayingStepRate = 0.1f;
    private float playingStep = 0;

    void Awake() {
        if (Instance != null && Instance != this) {
            Destroy(this);
        } else {
            Instance = this;
        }
    }

    void FixedUpdate() {
        localTime += Time.fixedDeltaTime;
        if (!Paused && localTime > playingStep) {
            playingStep = localTime + PlayingStepRate;
            SelectedPacketID += 1;
            
            SelectedPacket = DatabaseAccess.Instance.RecievedPackets.ContainsKey(SelectedPacketID) ? DatabaseAccess.Instance.RecievedPackets[SelectedPacketID] : null;
        }
    }
}
