using UnityEngine;
using UnityEngine.EventSystems;

public class DraggablePanel : MonoBehaviour, IBeginDragHandler, IDragHandler, IEndDragHandler
{
    private Vector2 originalLocalPointerPosition;
    private Vector2 originalPanelLocalPosition;
    private RectTransform panelRectTransform;
    private RectTransform parentRectTransform;
    
    // Threshold (in pixels) below which snapping occurs.
    public float snapThreshold = 20f;
    
    // Define a grid cell size for snapping (change these values as needed)
    public Vector2 snapGridSize = new Vector2(100f, 100f);

    void Awake()
    {
        panelRectTransform = transform as RectTransform;
        parentRectTransform = panelRectTransform.parent as RectTransform;
    }

    public void OnBeginDrag(PointerEventData data)
    {
        originalPanelLocalPosition = panelRectTransform.localPosition;
        RectTransformUtility.ScreenPointToLocalPointInRectangle(parentRectTransform, data.position, data.pressEventCamera, out originalLocalPointerPosition);
    }

    public void OnDrag(PointerEventData data)
    {
        if (panelRectTransform == null || parentRectTransform == null)
            return;

        Vector2 localPointerPosition;
        if (RectTransformUtility.ScreenPointToLocalPointInRectangle(parentRectTransform, data.position, data.pressEventCamera, out localPointerPosition))
        {
            Vector2 offsetToOriginal = localPointerPosition - originalLocalPointerPosition;
            panelRectTransform.localPosition = (Vector3)(originalPanelLocalPosition + offsetToOriginal);
        }
    }

    public void OnEndDrag(PointerEventData eventData)
    {
        // Snapping logic: Adjust the panel's position if it's close to a defined grid point.
        Vector3 pos = panelRectTransform.localPosition;
        float snappedX = Mathf.Round(pos.x / snapGridSize.x) * snapGridSize.x;
        float snappedY = Mathf.Round(pos.y / snapGridSize.y) * snapGridSize.y;
        Vector3 snappedPos = new Vector3(snappedX, snappedY, pos.z);
        
        // Apply snapping only if within the threshold distance
        if (Vector3.Distance(pos, snappedPos) < snapThreshold)
        {
            panelRectTransform.localPosition = snappedPos;
        }
    }
}
