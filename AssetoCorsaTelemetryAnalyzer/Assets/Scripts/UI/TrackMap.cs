using System.Collections.Generic;
using System.IO;
using UnityEngine;

public class TrackMap : MonoBehaviour {
    [Header("Refs")]
    public TrackMapLineRenderer lineRenderer;

    [Header("Map Options")]
    public int Width;
    public int Height;
    public Vector2 ManualOffset;

    [Space]
    [HideInInspector]
    public float[] TrackData = new float[]{
        0, 0,
        10, 10,
        20, 15,
        30, 20,
        40, 20,
        41, 40,
        43, 60,
        45, 80,
        50, 100,
        30, 120,
        15, 100,
        0, 100,
        0, 50,
    };

    [HideInInspector]
    public float scaleFactor;
    [HideInInspector]
    public float centerOffsetX;
    [HideInInspector]
    public float centerOffsetY;
    [HideInInspector]
    public float minX;
    [HideInInspector]
    public float minY;

    public void DrawTrackMap(float[] mapPoints) {
        minX = int.MaxValue;
        float maxX = int.MinValue;
        minY = int.MaxValue;
        float maxY = int.MinValue;
        for (int xIdx = 0; xIdx < mapPoints.Length; xIdx += 2) {
            int yIdx = xIdx + 1;
            float x = mapPoints[xIdx];
            float y = mapPoints[yIdx];

            minX = Mathf.Min(minX, x);
            maxX = Mathf.Max(maxX, x);
            minY = Mathf.Min(minY, y);
            maxY = Mathf.Max(maxY, y);
        }
        

        scaleFactor = Mathf.Min(Width/(maxX - minX), Height/(maxY - minY)) * 0.9f;
        bool boolScaleFavorX = Width/(maxX - minX) < Height/(maxY - minY);

        centerOffsetX = Width * 0.05f - Width/2;
        centerOffsetY = Height * 0.05f - Height/2;

        if (boolScaleFavorX) {
            centerOffsetY += (Height - ((maxY - minY) * scaleFactor))/2;
        } else {
            centerOffsetX += (Width - ((maxX - minX) * scaleFactor))/2;
        }

        List<Vector2> points = new List<Vector2>();

        for (int xIdx = 0; xIdx < mapPoints.Length; xIdx += 2) {
            int yIdx = xIdx + 1;
            float x = (mapPoints[xIdx] - minX) * scaleFactor + centerOffsetX + ManualOffset.x;
            float y = (mapPoints[yIdx] - minY) * scaleFactor + centerOffsetY + ManualOffset.y;

            points.Add(new Vector2(x, y));
        }

        lineRenderer.SetPoints(points);
    }

    void Update() {
        DrawTrackMap(TrackData);
    }
}
