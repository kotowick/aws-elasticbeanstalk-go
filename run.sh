
if [ -n "$VOLUME_PATH" ]; then
	echo "Sourcing any .env files in $VOLUME_PATH"
	source $VOLUME_PATH/*.env
fi

go-deploy -v upsert environment
