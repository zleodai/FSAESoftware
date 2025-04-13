using System.IO;
using UnityEngine;

public class TrackMap : MonoBehaviour {
    [Header("Refs")]
    public RenderTexture textureRender;

    [Header("Map Options")]
    public int Width;
    public int Height;
    public float LineWidth;

    private string logString;

    public void DrawTrackMap(float[] mapPoints) {
        float[,] trackMap = new float[Width, Height];

        float minX = int.MaxValue;
        float maxX = int.MinValue;
        float minY = int.MaxValue;
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
        

        float scaleFactor = Mathf.Min(Width/(maxX - minX), Height/(maxY - minY)) * 0.9f;
        bool boolScaleFavorX = Width/(maxX - minX) < Height/(maxY - minY);

        float centerOffsetX = Width * 0.05f;
        float centerOffsetY = Height * 0.05f;

        if (boolScaleFavorX) {
            centerOffsetY += scaleFactor/(Height/(maxY - minY)) * Height;
        } else {
            centerOffsetX += scaleFactor/(Width/(maxX - minX)) * Width;
        }

        float[] edges = new float[mapPoints.Length*2];

        for (int xIdx = 0; xIdx < mapPoints.Length; xIdx += 2) {
            int yIdx = xIdx + 1;
            float x = (mapPoints[xIdx] - minX) * scaleFactor + centerOffsetX;
            float y = (mapPoints[yIdx] - minY) * scaleFactor + centerOffsetY;

            float oldX;
            float oldY;
            
            if (xIdx == 0) {
                oldX = (mapPoints[mapPoints.Length - 2] - minX) * scaleFactor + centerOffsetX;
                oldY = (mapPoints[mapPoints.Length - 1] - minY) * scaleFactor + centerOffsetY;
            } else {
                oldX = (mapPoints[xIdx - 2] - minX) * scaleFactor + centerOffsetX;
                oldY = (mapPoints[yIdx - 2] - minY) * scaleFactor + centerOffsetY;
            }
            DrawLineOnMap((int)oldX, (int)oldY, (int)x, (int)y, trackMap);
        }

        Texture2D texture = new Texture2D (Width, Height);

        Color[] colorMap = new Color[Width * Height];
        for (int y = 0; y < Height; y++) {
            for (int x = 0; x < Width; x++) {
                colorMap [y * Width + x] = Color.Lerp(Color.black, Color.white, trackMap [x, y]);
            }
        }

        texture.SetPixels(colorMap);
        texture.Apply();
        Graphics.Blit(texture, textureRender);
    }

    private void DrawLineOnMap(int x1, int y1, int x2, int y2, float[,] map) {
        float x = x1 + 0.00001f;
        float y = y1 + 0.00001f;
        float dx = x2 - x1;
        float dy = y2 - y1;
        float d = 2 * dy - dx;
        float D = 0;
        float length = Mathf.Sqrt(dx * dx + dy * dy); 
        float sin = dx / length;
        float cos = dy / length;
        while (x <= x2) {
            if ( (int) x > 0 && (int) y - 1 > 0 && (int) x < map.GetLength(0) && (int) y + 1 < map.GetLength(1) ) { 
                map[(int) x, (int) y - 1] =  D + cos;
                map[(int) x, (int) y] =  D;
                map[(int) x, (int) y + 1] =  D - cos;
            }
            x = x + 1;
            if (d <= 0) {
                D += sin;
                d += 2 * dy;
            } else {
                D = D + sin - cos;
                d += 2 * (dy - dx);
                y++;
            }
        }
    }

    void Awake() {
        float[] sampleData = new float[]{
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
        DrawTrackMap(sampleData);
        File.WriteAllText("logdump.txt", logString);
    }
}
