using UnityEngine;

public class SizeVisualizer : MonoBehaviour
{
    public GraphController gc;
    private void Start()
    {
        RectTransform rt = GetComponent(typeof(RectTransform)) as RectTransform;
        rt.sizeDelta = new Vector2(gc.xgraphmax, gc.ygraphmax);
    }
}
