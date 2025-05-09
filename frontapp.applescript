--property lastApp : ""
property lastData : ""

on checkFrontApp()
	tell application "System Events"
		set currentApp to name of first application process whose frontmost is true
	end tell
	set currentData to ""
	
	if (currentApp = "Google Chrome") then
		tell application "Google Chrome"
			if it is running then
				set currentData to URL of active tab of front window
			end if
		end tell
	else if (currentApp = "Safari") then
		tell application "Safari"
			if it is running then
				if (count of windows) > 0 then
					set currentData to URL of current tab of window 1
				end if
			end if
		end tell
	end if
	
	set currentData to currentApp & "#" & currentData
	if currentData is not equal to lastData then
		set lastData to currentData
		log currentData
        set curlCommand to "curl -X POST -H \"Content-Type: text/plain\" --data \"" & currentData & "\" http://127.0.0.1:9482"
        -- log curlCommand
        try
            set serverResponse to do shell script curlCommand
        on error errorMessage number errorNumber
            log "errorMessage: " & errorMessage & " (errCode: " & errorNumber & ")"
        end try
		delay 0.5
	end if
	
	
end checkFrontApp

repeat
	checkFrontApp()
	delay 0.2
end repeat


