using System.Collections;
using UnityEngine;
using UnityEngine.UI;

namespace CarSpeedometerPro.Speedometer
{
    [ExecuteInEditMode]
    public class SpeedometerController : MonoBehaviour
    {
        [Space]
        [Header("----- SPEEDOMETER -----")]
        public Image speedometerSliderImage; // Image for the speedometer slider
        public Color speedometerSliderColor = Color.white; // Default color of the speedometer slider
        [Range(0f, 1f)]
        public float speedometerSliderValue = 0f; // Speedometer slider value

        [Space]
        [Header("Speedometer Text")]
        public Text currentSpeedText; // Text displaying the current speed
        public Text currentGearText; // Text displaying the current gear

        [Space]
        [Header("GAS SLIDER")]
        public Image gasSliderImage; // Image for the gas slider
        public Color gasSliderColor = Color.white; // Default color of the gas slider
        [Range(0f, 1f)]
        public float gasSliderValue = 0f; // Gas level value

        [Space]
        [Header("DAMAGE SLIDER")]
        public Image damageSliderImage; // Image for the damage slider
        public Color damageSliderColor = Color.white; // Default color of the damage slider
        [Range(0f, 1f)]
        public float damageSliderValue = 0f; // Damage level value

        [Space]
        [Header("Headlight States")]
        public GameObject offState;      // Represents the off state of headlights
        public GameObject lowBeamState;  // Represents the low beam state of headlights
        public GameObject highBeamState; // Represents the high beam state of headlights

        private int currentState = 0; // Tracks the current headlight state: 0 = Off, 1 = Low Beam, 2 = High Beam
        private GameObject[] states; // Array to store headlight state references

        [Space]
        [Header("HandBrake")]
        public Image handBrake; // Image representing the handbrake status
        private bool isHandBrakeActive = false; // Tracks whether the handbrake is engaged

        [Space]
        [Header("Blinking Settings")]
        public Image leftSignal;  // Image for the left turn signal
        public Image rightSignal; // Image for the right turn signal
        public float timeBlinking = 0.5f; // Blinking interval in seconds
        private Coroutine leftBlinkCoroutine; // Coroutine for left signal blinking
        private Coroutine rightBlinkCoroutine; // Coroutine for right signal blinking

        void OnValidate()
        {
            // Update UI elements when values are modified in the Inspector
            SetGasSliderColor();
            SetGasSliderProgress();
            SetDamageSliderColor();
            SetDamageSliderProgress();
            SetSpeedometerSliderColor();
            SetSpeedometerSliderProgress();
        }

        void SetGasSliderColor() => gasSliderImage.color = gasSliderColor;
        void SetGasSliderProgress() => gasSliderImage.fillAmount = Mathf.Lerp(0f, 1f, gasSliderValue);

        void SetDamageSliderColor() => damageSliderImage.color = damageSliderColor;
        void SetDamageSliderProgress() => damageSliderImage.fillAmount = Mathf.Lerp(0f, 1f, damageSliderValue);

        void SetSpeedometerSliderColor() => speedometerSliderImage.color = speedometerSliderColor;
        void SetSpeedometerSliderProgress() => speedometerSliderImage.fillAmount = Mathf.Lerp(0f, 0.75f, speedometerSliderValue);

        private void Start()
        {
            // Store headlight states in an array for easy switching
            states = new GameObject[] { offState, lowBeamState, highBeamState };

            // Ensure only the off state is active at startup
            SetState(0);
        }

        public void ToggleHeadlight()
        {
            // Cycle through headlight states
            currentState = (currentState + 1) % states.Length;
            SetState(currentState);
        }

        private void SetState(int index)
        {
            // Disable all states first
            foreach (GameObject state in states)
            {
                state.SetActive(false);
            }

            // Activate the selected state
            states[index].SetActive(true);
        }

        public void ToggleHandBrake()
        {
            // Toggle handbrake state and update the UI color accordingly
            isHandBrakeActive = !isHandBrakeActive;
            handBrake.color = isHandBrakeActive ? Color.red : Color.white;
        }

        public void ToggleLeftSignal()
        {
            if (leftBlinkCoroutine != null)
            {
                // Stop blinking and reset to default color
                StopCoroutine(leftBlinkCoroutine);
                leftSignal.color = Color.white;
                leftBlinkCoroutine = null;
            }
            else
            {
                // If right signal is active, stop it first
                if (rightBlinkCoroutine != null)
                    ToggleRightSignal();

                leftBlinkCoroutine = StartCoroutine(BlinkEffect(leftSignal));
            }
        }

        public void ToggleRightSignal()
        {
            if (rightBlinkCoroutine != null)
            {
                // Stop blinking and reset to default color
                StopCoroutine(rightBlinkCoroutine);
                rightSignal.color = Color.white;
                rightBlinkCoroutine = null;
            }
            else
            {
                // If left signal is active, stop it first
                if (leftBlinkCoroutine != null)
                    ToggleLeftSignal();

                rightBlinkCoroutine = StartCoroutine(BlinkEffect(rightSignal));
            }
        }

        private IEnumerator BlinkEffect(Image signal)
        {
            bool isYellow = true;

            while (true)
            {
                // Alternate between yellow and white to create a blinking effect
                signal.color = isYellow ? Color.yellow : Color.white;
                isYellow = !isYellow;
                yield return new WaitForSeconds(timeBlinking);
            }
        }
    }
}