if [ -n "$VOLUME_PATH" ]; then
	if ls $VOLUME_PATH/*.env > /dev/null 2>&1; then
		echo "Sourcing any .env files in $VOLUME_PATH"
		source $VOLUME_PATH/*.env
	fi
fi

go-deploy -v upsert
