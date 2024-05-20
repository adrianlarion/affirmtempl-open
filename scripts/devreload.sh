#!/bin/bash

CURRENT_WID=$(xdotool getwindowfocus)

WID=$(xdotool search --name "Google Chrome")
xdotool windowactivate $WID
sleep 0.8
xdotool key F5

xdotool windowactivate $CURRENT_WID
