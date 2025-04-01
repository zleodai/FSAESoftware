using UnityEngine;
using UnityEngine.InputSystem;

public class InputManager : MonoBehaviour
{
    void Start() {
        InputAction pausePlayAction = new InputAction("pausePlay", binding: "<Keyboard>/space");
        pausePlayAction.performed += ctx => onPausePlay();
        pausePlayAction.Enable();
    }

    void onPausePlay() {
        MainManager.Instance.Paused = !MainManager.Instance.Paused;
        
        if (MainManager.Instance.Paused) {
            Debug.Log("Paused");
        } else {
            Debug.Log("Playing");
        }
    }
}
