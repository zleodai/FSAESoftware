using UnityEngine;
using UnityEngine.EventSystems;

public class ResizablePanel : MonoBehaviour, IDragHandler, IPointerDownHandler
{
    // Reference to the panel you want to resize.
    public RectTransform target;
    private Vector2 originalLocalPointerPosition;
    private Vector2 originalSizeDelta;
    private Vector2 previousPointerPosition;
    
    public void OnPointerDown(PointerEventData data)
    {
        // Capture the initial state
        originalSizeDelta = target.sizeDelta;
        RectTransformUtility.ScreenPointToLocalPointInRectangle(target, data.position, data.pressEventCamera, out originalLocalPointerPosition);
        previousPointerPosition = data.position;
    }
    
    public void OnDrag(PointerEventData data)
    {
        if (target == null)
            return;

        // Calculate the delta movement in screen space
        Vector2 pointerDelta = data.position - previousPointerPosition;
        
        // Convert screen space delta to local space delta
        Vector2 localDelta;
        if (RectTransformUtility.ScreenPointToLocalPointInRectangle(target, data.position, data.pressEventCamera, out Vector2 localPoint) &&
            RectTransformUtility.ScreenPointToLocalPointInRectangle(target, previousPointerPosition, data.pressEventCamera, out Vector2 previousLocalPoint))
        {
            localDelta = localPoint - previousLocalPoint;
            
            // Update the size using the delta
            Vector2 newSizeDelta = target.sizeDelta + localDelta;
            
            // Clamp the minimum size
            newSizeDelta = new Vector2(Mathf.Max(newSizeDelta.x, 100), Mathf.Max(newSizeDelta.y, 100));
            target.sizeDelta = newSizeDelta;
        }
        
        // Update the previous position for the next frame
        previousPointerPosition = data.position;
    }
}
