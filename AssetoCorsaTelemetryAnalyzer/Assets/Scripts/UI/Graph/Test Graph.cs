using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class TestGraph : MonoBehaviour
{
    public LineRenderer lineGraph;

    void Start()
    {
        List<Vector2> points = new List<Vector2>();
        for (int x = 0; x <= 100; x++)
        {
            float y = Mathf.Sin(x * 0.1f) * 50f + 100f;
            points.Add(new Vector2(x * 5f, y));
        }

        lineGraph.SetPoints(points);
    }
}
