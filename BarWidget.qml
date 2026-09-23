import QtQuick
import Quickshell
import Quickshell.Io
import qs.Commons
import qs.Ui

BarWidget {
	id: root
	moduleName: "yairh.omaluah"
	implicitWidth: label.implicitWidth
	implicitHeight: barSize

	property string hebrewDate: ""

	Process {
		id: todayProc
		command: ["omaluah", "-today"]
		stdout: StdioCollector {
			waitForEnd: true 
			onStreamFinished: { 
				try {
					var day = JSON.parse(String(text || ""))
					root.hebrewDate = day ? day.hebrewDate : ""
				} catch (e) {}
			}
		}
	}

	Text {
		id: label
		text: root.hebrewDate || "שלום"
		anchors.verticalCenter: parent.verticalCenter
		anchors.horizontalCenter: parent.horizontalCenter
		color: root.bar ? root.bar.barForeground : Color.foreground
		font.pixelSize: Style.font.bodySmall
		font.family: root.bar ? root.bar.fontFamily : Style.font.family
	}

  function msToMidnight() {
    var next = new Date(); next.setHours(24, 0, 1, 0)
    return next.getTime() - Date.now()
  }

	Timer {
    id: refreshTimer
    interval: root.msToMidnight()
    repeat: true
    running: true
    triggeredOnStart: true
    onTriggered: { refreshTimer.interval = root.msToMidnight(); if (!todayProc.running) todayProc.running = true }
  }
}
