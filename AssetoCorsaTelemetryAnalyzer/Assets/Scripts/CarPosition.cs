using UnityEngine;

public class CarPosition : MonoBehaviour
{
    private RectTransform t;

    void Awake(){
        t = GetComponent<RectTransform>();
    }

    public void MovePosition(Vector2 p) {
        t.localPosition = p;
    }
}
