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
            speedometerController.throttleSliderValue = 0f;
            speedometerController.brakeSliderValue = 0f;
            speedometerController.currentSpeedText.text = "0";
        }
    }
}