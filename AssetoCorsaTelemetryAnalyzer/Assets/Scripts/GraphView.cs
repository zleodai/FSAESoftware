using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Networking;
using System.Collections;
using System.Collections.Generic;

public class GraphView : MonoBehaviour
{
    [Header("Graph Settings")]
    public RectTransform graphRect;    // The drawing area (usually the RectTransform of GraphArea)
    public float minY = 0f;            // Minimum value for the y-axis
    public float maxY = 100f;          // Maximum value for the y-axis
    public float xRange = 60f;         // Duration (in seconds) for the x-axis window

    [Header("Data Source Settings")]
    [Tooltip("Set this to true to fetch data from the PHP server; otherwise, random data is used.")]
    public bool useServerData = false;
    [Tooltip("URL for the PHP endpoint that returns graph data in JSON format.")]
    public string phpServerUrl = "http://yourserver.com/api/getgraphdata.php";

    // Reference to our custom line renderer component
    public SimpleUILineRenderer lineRenderer;

    // Data storage lists
    private List<float> indoorData = new List<float>();
    private List<float> outdoorData = new List<float>();
    private List<float> timeStamps = new List<float>();

    // Zoom factor for the x-axis
    private float zoomFactor = 1f;

    void Start()
    {
        if (useServerData)
        {
            StartCoroutine(GetDataFromServer());
        }
    }

    void Update()
    {
        // Use random data if not using the server.
        if (!useServerData)
        {
            float currentTime = Time.time;
            // Generate random placeholder values.
            AddDataPoint(currentTime, Random.Range(45f, 80f), Random.Range(30f, 90f));
            UpdateGraph();
        }
    }

    /// <summary>
    /// Coroutine for polling the PHP server for data.
    /// The PHP script should return JSON with fields "time", "indoor", and "outdoor".
    /// Example JSON: {"time": 123456.789, "indoor": 75.5, "outdoor": 68.3}
    /// </summary>
    private IEnumerator GetDataFromServer()
    {
        while (true)
        {
            UnityWebRequest www = UnityWebRequest.Get(phpServerUrl);
            yield return www.SendWebRequest();

#if UNITY_2020_1_OR_NEWER
            if (www.result != UnityWebRequest.Result.Success)
#else
            if (www.isNetworkError || www.isHttpError)
#endif
            {
                Debug.LogWarning("Failed to fetch data: " + www.error);
                float fallbackTime = Time.time;
                AddDataPoint(fallbackTime, Random.Range(45f, 80f), Random.Range(30f, 90f));
            }
            else
            {
                string jsonResponse = www.downloadHandler.text;
                DataPoint dataPoint = JsonUtility.FromJson<DataPoint>(jsonResponse);
                if (dataPoint != null)
                {
                    AddDataPoint(dataPoint.time, dataPoint.indoor, dataPoint.outdoor);
                }
                else
                {
                    Debug.LogWarning("Invalid JSON data received. Using fallback values.");
                    float fallbackTime = Time.time;
                    AddDataPoint(fallbackTime, Random.Range(45f, 80f), Random.Range(30f, 90f));
                }
            }

            UpdateGraph();
            yield return new WaitForSeconds(1.0f);
        }
    }

    /// <summary>
    /// Adds a new data point to the graph.
    /// </summary>
    public void AddDataPoint(float time, float indoorVal, float outdoorVal)
    {
        timeStamps.Add(time);
        // Here you can decide how to use multiple datasets.
        // For simplicity, we’re using one combined series here. You can extend this to support multiple lines.
        // E.g., you could choose to combine or toggle between indoorData and outdoorData.
        indoorData.Add(indoorVal);
        outdoorData.Add(outdoorVal);

        // Remove points that are outside the visible x-range.
        float cutoff = time - (xRange * zoomFactor);
        while (timeStamps.Count > 0 && timeStamps[0] < cutoff)
        {
            timeStamps.RemoveAt(0);
            indoorData.RemoveAt(0);
            outdoorData.RemoveAt(0);
        }
    }

    /// <summary>
    /// Updates the graph by converting data values to local UI coordinates.
    /// </summary>
    private void UpdateGraph()
    {
        if (timeStamps.Count == 0)
            return;

        float currentTime = timeStamps[timeStamps.Count - 1];
        float timeStart = currentTime - (xRange * zoomFactor);
        Vector2 size = graphRect.sizeDelta;

        // For demonstration, we’ll plot the indoor data line.
        // (You can add another renderer or draw multiple lines if needed.)
        List<Vector2> points = new List<Vector2>();

        for (int i = 0; i < timeStamps.Count; i++)
        {
            // Convert time to a normalized value (0 to 1)
            float normalizedX = Mathf.InverseLerp(timeStart, currentTime, timeStamps[i]);
            // Convert data value to a normalized y value (0 to 1)
            float normalizedY = Mathf.InverseLerp(minY, maxY, indoorData[i]);
            // Scale to actual pixel positions in the graph area.
            float xPos = normalizedX * size.x;
            float yPos = normalizedY * size.y;
            points.Add(new Vector2(xPos, yPos));
        }

        // Assign the calculated points to the custom line renderer.
        lineRenderer.Points = points;
        lineRenderer.SetVerticesDirty();

        //Debug.Log("First point: " + points[0] + " Last point: " + points[points.Count - 1]);
    }

    /// <summary>
    /// Optional: Call these methods from UI buttons to zoom in or out.
    /// </summary>
    public void ZoomIn()
    {
        zoomFactor *= 0.8f;
        if (zoomFactor < 0.1f) zoomFactor = 0.1f;
        UpdateGraph();
    }

    public void ZoomOut()
    {
        zoomFactor *= 1.25f;
        if (zoomFactor > 10f) zoomFactor = 10f;
        UpdateGraph();
    }
}

/// <summary>
/// Represents the data structure returned by your PHP backend.
/// Make sure the JSON from PHP matches these field names.
/// </summary>
[System.Serializable]
public class DataPoint
{
    public float time;
    public float indoor;
    public float outdoor;
}
