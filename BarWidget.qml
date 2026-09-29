import QtQuick
import Quickshell
import Quickshell.Io
import qs.Commons
import qs.Ui

BarWidget {
	id: root
	moduleName: "yairh.omaluah"
	implicitWidth: button.implicitWidth
	implicitHeight: barSize

	property string hebrewDate: ""
	property string holiday: ""

	function getHolidays(list) {
		if (!Array.isArray(list) || list.length === 0) { return }
		var titles = []
		for (var i = 0; i < list.length; i++) 
		if (list[i].title) {
			titles.push(String(list[i].title)) 
			titles.push(String(list[i].hebrewTitle))
		}
		return titles.join(" | ")
	}

	function refresh() {
		if (!todayProc.running) todayProc.running = true
	}

	Process {
		id: todayProc
		command: ["omaluah", "-today"]
		stdout: StdioCollector {
			waitForEnd: true 
			onStreamFinished: { 
				try {
					var day = JSON.parse(String(text || ""))
					root.hebrewDate = day ? day.hebrewDate : ""
					root.holiday = day ? root.getHolidays(day.holidays) : ""
				} catch (e) {}
			}
		}
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
    onTriggered: { refreshTimer.interval = root.msToMidnight(); root.refresh() }
  }

  // ---- Calendar popup. Shape contract for shell.summon/hide/toggle
  //      routing: Bar.findPanelWidget requires open/close/opened on the
  //      bar-widget root.
  readonly property bool opened: panelLoader.item ? panelLoader.item.opened === true : false

  function open() {
    if (panelLoader.item) panelLoader.item.open()
  }

  function close() {
    if (panelLoader.item) panelLoader.item.close()
  }

  function togglePanel() {
    if (panelLoader.item) panelLoader.item.toggle()
  }

  function injectPanel() {
    var target = panelLoader.item
    if (!target) return
    if ("bar" in target) target.bar = root.bar
    if ("settings" in target) target.settings = root.settings
    if ("anchorItem" in target) target.anchorItem = button
    if ("hostWidget" in target) target.hostWidget = root
  }

  Loader {
    id: panelLoader
    active: true
    source: Qt.resolvedUrl("Panel.qml")
    visible: false
    onLoaded: {
      root.injectPanel()
      Qt.callLater(root.injectPanel)
    }
  }

  onBarChanged: injectPanel()
  onSettingsChanged: injectPanel()

  WidgetButton {
    id: button
    anchors.fill: parent
    bar: root.bar
		text: root.hebrewDate || "שלום"
		fontSize: Style.font.bodySmall
    horizontalMargin: 8.75
    verticalPadding: 8.75

    onPressed: function(b) {
      root.togglePanel()
    }
	}
}
