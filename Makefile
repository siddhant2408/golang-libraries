include Makefile-common.mk

ALL_DIR=ls -d -1 ../*/ | grep -v golang-libraries
ALL_COMMAND=$(ALL_DIR) | xargs -n 1 -I {} $(1)
ALL_RUN=$(call ALL_COMMAND,sh -c 'echo {} && cd {} && $(1)')
.PHONY: all-run
all-run:
	$(eval COMMAND?=ls)
	$(call ALL_RUN,$(COMMAND))

.PHONY: all-build
all-build:
	$(call ALL_RUN,make build)

.PHONY: all-test
all-test:
	$(call ALL_RUN,make test)

.PHONY: all-lint
all-lint:
	$(call ALL_RUN,make lint)

.PHONY: all-golangci-lint
all-golangci-lint: install-golangci-lint
	$(call ALL_RUN,make golangci-lint)

.PHONY: all-lint-rules
all-lint-rules:
	$(call ALL_RUN,make lint-rules)

.PHONY: all-clean
all-clean:
	$(call ALL_RUN,make clean)

.PHONY: all-mod-update
all-mod-update: all-copy-common
	$(call ALL_RUN,make mod-update)

.PHONY: all-mod-update-golang-libraries
all-mod-update-golang-libraries: all-copy-common
	$(call ALL_RUN,make mod-update-golang-libraries)

.PHONY: all-mod-tidy
all-mod-tidy:
	$(call ALL_RUN,make mod-tidy)

.PHONY: all-mod-edit-replace-golang-libraries-local
all-mod-edit-replace-golang-libraries-local:
	$(call ALL_RUN,make mod-edit-replace-golang-libraries-local)

.PHONY: all-mod-edit-dropreplace-golang-libraries
all-mod-edit-dropreplace-golang-libraries:
	$(call ALL_RUN,make mod-edit-dropreplace-golang-libraries)

.PHONY: all-git-latest-release
all-git-latest-release:
	$(call ALL_RUN,make git-latest-release)

.PHONY: all-copy-common
all-copy-common:
	$(call ALL_COMMAND,cp -r Makefile-common.mk Makefile-project.mk .gitignore .github .dockerignore .golangci.yml {})

ifeq ($(CI),true)

ci::
	$(call CI_LOG_GROUP_START,test)
	$(MAKE) test
	$(call CI_LOG_GROUP_END)

	$(call CI_LOG_GROUP_START,lint)
	$(MAKE) lint
	$(call CI_LOG_GROUP_END)

ci-services:: ci-service-clickhouse ci-service-mongodb ci-service-mysql ci-service-rabbitmq ci-service-redis

endif # CI end
