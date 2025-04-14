using UnityEngine;
using CarSpeedometerPro.Speedometer;

namespace CarSpeedometerPro.UI
{
    public class SpeedometerManager : MonoBehaviour
    {
        private SpeedometerController speedometerController;

        void Awake()
        {
            speedometerController = FindObjectOfType<SpeedometerController>();
        }

        public void Start()
        {
            // Initialize SPEEDOMETER values
            speedometerController.speedometerSliderValue = 0.3f;
            speedometerController.currentGearText.text = "N";
            speedometerController.currentSpeedText.text = "0";

            // Initialize GAS value
            speedometerController.gasSliderValue = 1f;

            // Initialize DAMAGE value
            speedometerController.damageSliderValue = 1f;
        }

        public void ToggleLeftSignal()
        {
            speedometerController.ToggleLeftSignal();
        }

        public void ToggleRightSignal()
        {
            speedometerController.ToggleRightSignal();
        }

        public void ToggleHandBrake()
        {
            speedometerController.ToggleHandBrake();
        }

        public void ToggleHeadlights()
        {
            speedometerController.ToggleHeadlight();
        }
    }
}