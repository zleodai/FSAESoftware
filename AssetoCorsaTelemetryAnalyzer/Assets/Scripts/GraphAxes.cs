using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

public class GraphAxes : MonoBehaviour
{
    [Header("Graph Container (UI Panel)")]
    public RectTransform graphContainer;

    [Header("Graph Margins (pixels)")]
    public float marginLeft = 100f;
    public float marginRight = 20f;
    public float marginTop = 20f;
    public float marginBottom = 100f;

    [Header("Axis Settings")]
    public float xMin = 0f;
    public float xMax;
    public float yMin = 0f;
    public float yMax;
    public int xAxisDivisions = 10;
    public int yAxisDivisions = 10;

    [Header("UI Settings")]
    public Sprite lineSprite;
    public Font labelFont;
    public int fontSize = 16;
    public string xAxisLabel = "X Axis";
    public string yAxisLabel = "Y Axis";
    public float labelSpacing = 25f;

    private List<GameObject> labelObjects = new List<GameObject>();
    private List<GameObject> lineObjects = new List<GameObject>();
    private float graphWidth;
    private float graphHeight;
    private GraphController graphController;

    void Awake()
    {
        graphController = GetComponentInParent<GraphController>();
        if (graphController == null)
        {
            Debug.LogError("GraphController not found!");
            return;
        }

        RectTransform rt = GetComponent<RectTransform>();
        if (rt != null)
        {
            // Set the full size including margins
            graphWidth = graphController.xgraphmax + marginLeft + marginRight;
            graphHeight = graphController.ygraphmax + marginTop + marginBottom;
            rt.sizeDelta = new Vector2(graphWidth, graphHeight);
            
            xMax = graphController.xmax;
            yMax = graphController.ymax;
        }
    }

    void Start()
    {
        CreateAxes();
    }

    public void CreateAxes()
    {
        foreach (GameObject obj in labelObjects) Destroy(obj);
        labelObjects.Clear();
        foreach (GameObject obj in lineObjects) Destroy(obj);
        lineObjects.Clear();

        float xAxisLength = graphController.xgraphmax;
        float yAxisLength = graphController.ygraphmax;

        // Draw X Axis at the bottom margin
        Vector2 xStart = new Vector2(marginLeft, marginBottom);
        Vector2 xEnd = new Vector2(marginLeft + xAxisLength, marginBottom);
        GameObject xAxis = CreateLine("X Axis", xStart, xEnd);
        lineObjects.Add(xAxis);

        // Draw Y Axis at the left margin
        Vector2 yStart = new Vector2(marginLeft, marginBottom);
        Vector2 yEnd = new Vector2(marginLeft, marginBottom + yAxisLength);
        GameObject yAxis = CreateLine("Y Axis", yStart, yEnd);
        lineObjects.Add(yAxis);

        // X-Axis tick marks and labels
        float tickHeight = 5f;
        for (int i = 0; i <= xAxisDivisions; i++)
        {
            float t = i / (float)xAxisDivisions;
            float xPos = marginLeft + (t * xAxisLength);
            float dataValue = t * xMax;

            // Draw tick
            GameObject tick = CreateLine("X Tick " + i,
                new Vector2(xPos, marginBottom),
                new Vector2(xPos, marginBottom - tickHeight));
            lineObjects.Add(tick);

            // Create label
            Vector2 labelPos = new Vector2(xPos, marginBottom - tickHeight - labelSpacing);
            GameObject label = CreateText(dataValue.ToString("0"), labelPos, TextAnchor.UpperCenter);
            labelObjects.Add(label);
        }

        // Y-Axis tick marks and labels
        float tickWidth = 5f;
        for (int i = 0; i <= yAxisDivisions; i++)
        {
            float t = i / (float)yAxisDivisions;
            float yPos = marginBottom + (t * yAxisLength);
            float dataValue = t * yMax;

            // Draw tick
            GameObject tick = CreateLine("Y Tick " + i,
                new Vector2(marginLeft, yPos),
                new Vector2(marginLeft - tickWidth, yPos));
            lineObjects.Add(tick);

            // Create label - using same spacing as x-axis
            Vector2 labelPos = new Vector2(marginLeft - tickWidth - (labelSpacing * 2), yPos);
            GameObject label = CreateText(dataValue.ToString("0"), labelPos, TextAnchor.MiddleRight);
            labelObjects.Add(label);
        }

        // Add axis labels
        GameObject xTitle = CreateText(xAxisLabel,
            new Vector2(marginLeft + xAxisLength / 2f, marginBottom / 2f),
            TextAnchor.MiddleCenter);
        labelObjects.Add(xTitle);

        GameObject yTitle = CreateText(yAxisLabel,
            new Vector2(marginLeft / 2f - (labelSpacing * 1.5f), marginBottom + yAxisLength / 2f),
            TextAnchor.MiddleCenter);
        yTitle.GetComponent<RectTransform>().localEulerAngles = new Vector3(0, 0, 90);
        labelObjects.Add(yTitle);
    }

    private GameObject CreateText(string textString, Vector2 anchoredPos, TextAnchor alignment)
    {
        GameObject textGO = new GameObject("Label", typeof(Text));
        textGO.transform.SetParent(transform, false);
        Text text = textGO.GetComponent<Text>();
        text.text = textString;
        text.font = labelFont;
        text.fontSize = fontSize;
        text.alignment = alignment;
        text.color = Color.black;

        RectTransform rt = textGO.GetComponent<RectTransform>();
        rt.anchorMin = Vector2.zero;
        rt.anchorMax = Vector2.zero;
        rt.pivot = new Vector2(0.5f, 0.5f);
        rt.anchoredPosition = anchoredPos;
        rt.sizeDelta = new Vector2(60, 30);
        return textGO;
    }

    private GameObject CreateLine(string name, Vector2 startPoint, Vector2 endPoint)
    {
        GameObject line = new GameObject(name, typeof(Image));
        line.transform.SetParent(transform, false);
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
