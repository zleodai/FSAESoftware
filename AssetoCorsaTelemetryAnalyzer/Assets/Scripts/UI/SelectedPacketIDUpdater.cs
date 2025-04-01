using TMPro;
using UnityEngine;

public class SelectedPacketIDUpdater : MonoBehaviour
{    
    TextMeshProUGUI text;

    void Start() {
        text = gameObject.GetComponent<TextMeshProUGUI>();
    }
    
    void Update() {
        if (text != null) {
            text.text = string.Format("SelectedPacketID: {0}", MainManager.Instance.SelectedPacketID);
        }
    }
}
