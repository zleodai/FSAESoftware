using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class LineRenderer : Graphic
{
    public GraphController gc;
    public float marginLeft = 100f;
    public float marginBottom = 100f;

    public List<Vector2> Points = new List<Vector2>();
    public float thickness = 2f;

    // Max values of the data
    public float xmax = 100f;
    public float ymax = 100f;

    // Dimensions in pixels that the graph should occupy
    public float xgraphmax = 300f;
    public float ygraphmax = 200f;

    protected override void Start()
    {
        if (gc != null)
        {
            RectTransform rt = GetComponent<RectTransform>();
            // Set the full size including margins
            rt.sizeDelta = new Vector2(gc.xgraphmax + marginLeft + 20f, gc.ygraphmax + marginBottom + 20f);
            xmax = gc.xmax;
            ymax = gc.ymax;
            xgraphmax = gc.xgraphmax;
            ygraphmax = gc.ygraphmax;
        }
    }

    protected override void OnPopulateMesh(VertexHelper vh)
    {
        vh.Clear();

        if (Points == null || Points.Count < 2)
            return;

        // Convert data points into graph-space positions
        List<Vector2> scaledPoints = new List<Vector2>();
        foreach (Vector2 point in Points)
        {
            // Scale the data points to fit the graph area and offset by margins
            float x = marginLeft + (point.x / xmax * xgraphmax);
            float y = marginBottom + (point.y / ymax * ygraphmax);
            scaledPoints.Add(new Vector2(x, y));
        }

        for (int i = 0; i < scaledPoints.Count - 1; i++)
        {
            Vector2 p1 = scaledPoints[i];
            Vector2 p2 = scaledPoints[i + 1];
            DrawLine(vh, p1, p2, thickness);
        }
    }

    void DrawLine(VertexHelper vh, Vector2 p1, Vector2 p2, float thickness)
    {
        Vector2 direction = (p2 - p1).normalized;
        Vector2 normal = new Vector2(-direction.y, direction.x) * (thickness / 2f);

        Vector2 v0 = p1 - normal;
        Vector2 v1 = p1 + normal;
        Vector2 v2 = p2 + normal;
        Vector2 v3 = p2 - normal;

        int startIndex = vh.currentVertCount;

        vh.AddVert(v0, color, Vector2.zero);
        vh.AddVert(v1, color, Vector2.zero);
        vh.AddVert(v2, color, Vector2.zero);
        vh.AddVert(v3, color, Vector2.zero);

        vh.AddTriangle(startIndex, startIndex + 1, startIndex + 2);
        vh.AddTriangle(startIndex + 2, startIndex + 3, startIndex);
    }

    public void SetPoints(List<Vector2> newPoints)
    {
        Points = newPoints;
        SetVerticesDirty();
    }
}
