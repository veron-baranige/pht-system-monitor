SERVICE_NAME=springboot-app-monitor

install:
	@mkdir -p ~/.config/systemd/user/
	@sudo mkdir -p /usr/local/bin/${SERVICE_NAME}
	@sudo cp -r ./ /usr/local/bin/${SERVICE_NAME}/
	@sudo chmod +x /usr/local/bin/${SERVICE_NAME}/bin
	@sudo chmod +x /usr/local/bin/${SERVICE_NAME}/start.sh
	@cp ./config/monitor.service ~/.config/systemd/user/${SERVICE_NAME}.service
	@systemctl --user daemon-reload
	@systemctl --user enable ${SERVICE_NAME}
	@systemctl --user start ${SERVICE_NAME}

uninstall:
	@systemctl --user stop ${SERVICE_NAME} || true
	@systemctl --user disable ${SERVICE_NAME} || true
	@rm -f ~/.config/systemd/user/${SERVICE_NAME}.service
	@systemctl --user daemon-reload
	@sudo rm -rf /usr/local/bin/${SERVICE_NAME}