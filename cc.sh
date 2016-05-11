#!/bin/bash

if [ -x ~/bin/color-cat ] ; then
	~/bin/color-cat -c green
else	
	cat
fi

