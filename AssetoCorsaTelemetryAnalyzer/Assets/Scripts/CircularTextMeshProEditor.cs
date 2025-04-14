using UnityEngine;
using UnityEditor;

[CustomEditor (typeof(CircularTextMeshPro))]
public class CircularTextMeshProEditor : Editor {
    public override void OnInspectorGUI() {
        CircularTextMeshPro circlethingy = (CircularTextMeshPro)target;

        if (DrawDefaultInspector()) {
            circlethingy.OnCurvePropertyChanged();
        }
    }
}
