using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class LineRenderer : Graphic
{
    public GraphController gc;

    public List<Vector2> Points = new List<Vector2>();
    public float thickness = 2f;

    // Max values of the data
    public int xmax = 100;
    public int ymax = 100;

    // Dimensions in pixels that the graph should occupy
    public float xgraphmax;
    public float ygraphmax;

    public int xmin = 0;
    public int ymin = 0;

    protected override void Start()
    {
        RectTransform rt = GetComponent(typeof(RectTransform)) as RectTransform;
        rt.sizeDelta = new Vector2(gc.xgraphmax, gc.ygraphmax);

        xmax = gc.xmax;
        ymax = gc.ymax;

        xmin = gc.xmin;
        ymin = gc.ymin;
    }

    protected override void OnPopulateMesh(VertexHelper vh)
    {
        vh.Clear();
        xgraphmax = rectTransform.rect.width;
        ygraphmax = rectTransform.rect.height;

        if (Points == null || Points.Count < 2)
            return;

        // Convert data points into graph-space positions
        List<Vector2> scaledPoints = new List<Vector2>();
        float xRange = Mathf.Max(1, xmax - xmin);
        float yRange = Mathf.Max(1, ymax - ymin);

        foreach (Vector2 point in Points)
        {
            float xNorm = Mathf.Clamp01((point.x - xmin) / xRange);
            float yNorm = Mathf.Clamp01((point.y - ymin) / yRange);

            float x = xNorm * xgraphmax;
            float y = yNorm * ygraphmax;

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
