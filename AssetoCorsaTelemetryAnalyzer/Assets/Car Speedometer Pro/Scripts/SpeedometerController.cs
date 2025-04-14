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
        public Image throttleSliderImage; 
        public Color throttleSliderColor = Color.white; 
        [Range(0f, 1f)]
        public float throttleSliderValue = 0f; 

        [Space]
        public Image brakeSliderImage; 
        public Color brakeSliderColor = Color.white; 
        [Range(0f, 1f)]
        public float brakeSliderValue = 0f; 

        [Space]
        [Header("Speedometer Text")]
        public Text currentSpeedText; // Text displaying the current speed

        void Update()
        {
            SetThrottleSliderColor();
            SetThrottleSliderProgress();
            SetBrakeSliderColor();
            SetBrakeSliderProgress();
        }

        void SetThrottleSliderColor() => throttleSliderImage.color = throttleSliderColor;
        void SetThrottleSliderProgress() => throttleSliderImage.fillAmount = Mathf.Lerp(0f, 0.75f, throttleSliderValue);
        void SetBrakeSliderColor() => brakeSliderImage.color = brakeSliderColor;
        void SetBrakeSliderProgress() => brakeSliderImage.fillAmount = Mathf.Lerp(0f, 0.75f, brakeSliderValue);
    }
}