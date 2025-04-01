using UnityEngine;
using System.Collections.Generic;

public class CarPosition : MonoBehaviour
{
    private long lastSelectedPacketID = -404;

    void Update() {
        long selectedPacketID = MainManager.Instance.SelectedPacketID;

        if (lastSelectedPacketID != selectedPacketID) {
            if (DatabaseAccess.Instance.RecievedPackets.ContainsKey(selectedPacketID)) {
                DatabaseAccess.Packet packet = DatabaseAccess.Instance.RecievedPackets[selectedPacketID];

                Vector3 position = new Vector3(packet.TelemetryData.WorldPositionX, packet.TelemetryData.WorldPositionY, packet.TelemetryData.WorldPositionZ);
                gameObject.transform.localPosition = position;
            } else {
                DatabaseAccess.Instance.QueryPackets(MainManager.Instance.SelectedPacketID, MainManager.Instance.SelectedPacketID + 100);
            }
            
            lastSelectedPacketID = selectedPacketID;
        }
    }
}
