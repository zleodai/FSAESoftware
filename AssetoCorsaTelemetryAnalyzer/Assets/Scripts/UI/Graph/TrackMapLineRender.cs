using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class TrackMapLineRenderer : Graphic
{
    public List<Vector2> Points = new List<Vector2>();
    public float thickness = 2f;

    protected override void OnPopulateMesh(VertexHelper vh)
    {
        vh.Clear();

        if (Points == null || Points.Count < 2)
            return;

        for (int i = 0; i < Points.Count - 1; i++)
        {
            Vector2 p1 = Points[i];
            Vector2 p2 = Points[i + 1];
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
