using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class GraphAxes : MonoBehaviour
{
    [Header("Graph Container (UI Panel)")]
    public RectTransform graphContainer;

    [Header("Graph Margins (pixels)")]
    public float marginLeft = 0f;  // Increased left margin for more label space
    public float marginRight = 0f;
    public float marginTop = 0f;
    public float marginBottom = 0f; // Increased bottom margin for larger labels

    [Header("Axis Settings")]
    public float xMin = 0f;
    public float xMax = 10f;
    public float yMin = 0f;
    public float yMax = 10f;
    public int xAxisDivisions = 10;
    public int yAxisDivisions = 10;

    [Header("UI Settings")]
    public Sprite lineSprite;
    public Font labelFont;
    public int fontSize = 16;  // Made font size configurable
    public string xAxisLabel = "X Axis";  // Made axis labels configurable
    public string yAxisLabel = "Y Axis";
    public float labelSpacing = 25f;  // Made label spacing configurable

    private List<GameObject> labelObjects = new List<GameObject>();
    private List<GameObject> lineObjects = new List<GameObject>();

    void Awake()
    {
        // Get the GraphController component from one of the parent GameObjects.
        GraphController gc = GetComponentInParent<GraphController>();
        //UNCOMMENT ^

        //Get the RectTransform attached to this GameObject.
        RectTransform rt = GetComponent(typeof(RectTransform)) as RectTransform;
        //UNCOMMENT ^

        if (gc != null && rt != null)
        {
            // Set the sizeDelta (width and height) to the values from GraphController.
            rt.sizeDelta = new Vector2(gc.xgraphmax, gc.ygraphmax);
        }
        else
        {
            Debug.LogWarning("GraphController or RectTransform not found.");
        }
        if (graphContainer != null)
        {
            // Make the container stretch with the panel
            graphContainer.anchorMin = Vector2.zero;
            graphContainer.anchorMax = Vector2.one;
            graphContainer.pivot = Vector2.zero;
            graphContainer.anchoredPosition = Vector2.zero;
            graphContainer.sizeDelta = Vector2.zero;
        }
    }

    void Start()
    {
        CreateAxes();
    }

    public void CreateAxes()
    {
        // Clear old elements
        foreach (GameObject obj in labelObjects) Destroy(obj);
        labelObjects.Clear();
        foreach (GameObject obj in lineObjects) Destroy(obj);
        lineObjects.Clear();

        float totalWidth = graphContainer.rect.width;
        float totalHeight = graphContainer.rect.height;

        float xAxisLength = totalWidth - marginLeft - marginRight;
        float yAxisLength = totalHeight - marginTop - marginBottom;

        if (xAxisLength < 0 || yAxisLength < 0)
        {
            Debug.LogWarning("Margins are too large for the panel size.");
            return;
        }

        Vector2 origin = new Vector2(marginLeft, marginBottom);

        // Draw X Axis
        GameObject xAxis = CreateLine("X Axis", origin, new Vector2(marginLeft + xAxisLength, marginBottom));
        lineObjects.Add(xAxis);

        // X-Axis tick marks and labels
        float tickHeight = 5f;
        float labelOffset = labelSpacing;
        for (int i = 0; i <= xAxisDivisions; i++)
        {
            float t = i / (float)xAxisDivisions;
            float xPos = Mathf.Lerp(origin.x, origin.x + xAxisLength, t);

            // Tick mark pointing down
            GameObject tick = CreateLine("X Tick " + i,
                new Vector2(xPos, marginBottom),
                new Vector2(xPos, marginBottom - tickHeight));
            lineObjects.Add(tick);

            // Label below the tick
            float currentVal = Mathf.Lerp(xMin, xMax, t);
            Vector2 labelPos = new Vector2(xPos, marginBottom - tickHeight - labelOffset);
            GameObject label = CreateText(currentVal.ToString("0.##"), labelPos, TextAnchor.UpperCenter);
            labelObjects.Add(label);
        }

        // Draw Y Axis
        GameObject yAxis = CreateLine("Y Axis", origin, new Vector2(marginLeft, marginBottom + yAxisLength));
        lineObjects.Add(yAxis);

        // Y-Axis tick marks and labels
        float tickWidth = 5f;
        labelOffset = labelSpacing;
        for (int i = 0; i <= yAxisDivisions; i++)
        {
            float t = i / (float)yAxisDivisions;
            float yPos = Mathf.Lerp(origin.y, origin.y + yAxisLength, t);

            // Tick mark pointing left
            GameObject tick = CreateLine("Y Tick " + i,
                new Vector2(marginLeft, yPos),
                new Vector2(marginLeft - tickWidth, yPos));
            lineObjects.Add(tick);

            // Label to the left of tick, moved further left
            float currentVal = Mathf.Lerp(yMin, yMax, t);
            Vector2 labelPos = new Vector2(marginLeft - tickWidth - labelOffset - 15f, yPos);  // Added extra offset
            GameObject label = CreateText(currentVal.ToString("0.##"), labelPos, TextAnchor.MiddleRight);
            labelObjects.Add(label);
        }

        // Axis titles
        GameObject xTitle = CreateText(xAxisLabel,
            new Vector2(marginLeft + xAxisLength / 2f, marginBottom / 2f),
            TextAnchor.MiddleCenter);
        labelObjects.Add(xTitle);

        GameObject yTitle = CreateText(yAxisLabel,
            new Vector2(marginLeft / 2f - 15f, marginBottom + yAxisLength / 2f),  // Moved Y axis title more left
            TextAnchor.MiddleCenter);
        yTitle.GetComponent<RectTransform>().localEulerAngles = new Vector3(0, 0, 90);
        labelObjects.Add(yTitle);
    }

    private GameObject CreateText(string textString, Vector2 anchoredPos, TextAnchor alignment)
    {
        GameObject textGO = new GameObject("Label", typeof(Text));
        textGO.transform.SetParent(graphContainer, false);
        Text text = textGO.GetComponent<Text>();
        text.text = textString;
        text.font = labelFont;
        text.fontSize = fontSize;  // Use the configurable font size
        text.alignment = alignment;
        text.color = Color.black;

        RectTransform rt = textGO.GetComponent<RectTransform>();
        rt.anchorMin = Vector2.zero;
        rt.anchorMax = Vector2.zero;
        rt.pivot = new Vector2(0.5f, 0.5f);
        rt.anchoredPosition = anchoredPos;
        rt.sizeDelta = new Vector2(60, 30);  // Increased text box size
        return textGO;
    }

    private GameObject CreateLine(string name, Vector2 startPoint, Vector2 endPoint)
    {
        GameObject line = new GameObject(name, typeof(Image));
        line.transform.SetParent(graphContainer, false);
        Image image = line.GetComponent<Image>();
        image.sprite = lineSprite;
        image.color = Color.black;

        RectTransform rt = line.GetComponent<RectTransform>();
        rt.anchorMin = Vector2.zero;
        rt.anchorMax = Vector2.zero;
        rt.pivot = new Vector2(0, 0.5f);
        
        Vector2 direction = (endPoint - startPoint).normalized;
        float distance = Vector2.Distance(startPoint, endPoint);
        
        rt.sizeDelta = new Vector2(distance, 2f);
        rt.anchoredPosition = startPoint;
        
        float angle = Mathf.Atan2(direction.y, direction.x) * Mathf.Rad2Deg;
        rt.localEulerAngles = new Vector3(0, 0, angle);

        return line;
    }
}
