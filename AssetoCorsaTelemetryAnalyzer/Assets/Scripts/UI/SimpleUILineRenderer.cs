using UnityEngine;
using UnityEngine.UI;
using System.Collections.Generic;

// Custom UI component for drawing a line from a set of points.
public class SimpleUILineRenderer : MaskableGraphic
{
    // List of points (in local UI coordinates) to connect
    public List<Vector2> Points = new List<Vector2>();

    // Thickness of the line in pixels
    public float Thickness = 2f;

    /// <summary>
    /// Generates the mesh based on the provided points.
    /// </summary>
    protected override void OnPopulateMesh(VertexHelper vh)
    {
        vh.Clear();
        Debug.Log("OnPopulateMesh called. Points count: " + (Points == null ? 0 : Points.Count));


        // Need at least two points to draw a line.
        if (Points == null || Points.Count < 2)
            return;

        // Iterate through the points and create a quad for each segment.
        for (int i = 0; i < Points.Count - 1; i++)
        {
            Vector2 start = Points[i];
            Vector2 end = Points[i + 1];

            // Calculate the direction vector and the perpendicular (normal) for thickness.
            Vector2 direction = (end - start).normalized;
            Vector2 normal = new Vector2(-direction.y, direction.x);
            Vector2 offset = normal * (Thickness * 0.5f);

            // Create four vertices for the quad representing the line segment.
            UIVertex vert = UIVertex.simpleVert;
            vert.color = color;

            // Lower left of the quad.
            vert.position = start - offset;
            vh.AddVert(vert);

            // Upper left.
            vert.position = start + offset;
            vh.AddVert(vert);

            // Upper right.
            vert.position = end + offset;
            vh.AddVert(vert);

            // Lower right.
            vert.position = end - offset;
            vh.AddVert(vert);

            // Calculate the starting index for this segment.
            int idx = i * 4;
            // Add two triangles to form the quad.
            vh.AddTriangle(idx, idx + 1, idx + 2);
            vh.AddTriangle(idx, idx + 2, idx + 3);
        }
    }
}
