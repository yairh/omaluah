import QtQuick
import Quickshell
import Quickshell.Io
import qs.Commons
import qs.Ui

Panel {
	id: root
	moduleName: "yairh.omaluah"

	property var anchorItem: null
  property var hostWidget: null
  readonly property var barIdentity: hostWidget || root

	function refresh() { if (hostWidget && hostWidget.refresh) hostWidget.refresh() }

  // Summoning by hotkey moves no pointer, so a hover the bar was still
  // holding must not keep the center indicators revealed behind the panel.
  function setCenterHoverRevealSuppressed(value) {
    if (root.bar && typeof root.bar.setCenterHoverRevealSuppressed === "function")
      root.bar.setCenterHoverRevealSuppressed(value)
    else if (root.bar && "centerHoverRevealSuppressed" in root.bar)
      root.bar.centerHoverRevealSuppressed = value
  }
	
  function open() {
    refresh()
    root.controller.show()
    // Set after showing, not before: showing hands the popout coordinator
    // over, which closes whichever panel was open, and that close clears the
    // shared flag. Deferring means the panel taking over always wins, while
    // a handoff to a panel that does not manage the flag still leaves it
    // cleared rather than stuck on.
    Qt.callLater(function() {
      if (root.opened) setCenterHoverRevealSuppressed(true)
    })
  }

  function close() {
    setCenterHoverRevealSuppressed(false)
    // Dismissing the panel mid-edit would otherwise leave the inputs up,
    // waiting behind a closed popup for the next time it opens.
    root.controller.hide()
  }

  function toggle() {
    if (root.opened) root.close()
    else root.open()
  }

  KeyboardPanel {
    id: panel
    anchorItem: root.anchorItem
    owner: root.barIdentity
    bar: root.bar
    open: root.opened
    focusTarget: keyCatcher
    contentWidth: panel.fittedContentWidth(Style.space(360))
    contentHeight: panel.fittedContentHeight(body.implicitHeight)

			PanelKeyCatcher {
				id: keyCatcher
				anchors.fill: parent

				Column {
					id: body
					anchors.left: parent.left
					anchors.right: parent.right
					anchors.top: parent.top
					anchors.verticalCenter: parent.verticalCenter
					spacing: Style.space(10)

					Text {
						id: label
						text: hostWidget.hebrewDate || "שלום"
						width: parent.width
						horizontalAlignment: Text.AlignHCenter
						color: root.bar ? root.bar.barForeground : Color.foreground
						font.pixelSize: Style.font.bodySmall
						font.family: root.bar ? root.bar.fontFamily : Style.font.family
					}
				}
			}
		}
}
