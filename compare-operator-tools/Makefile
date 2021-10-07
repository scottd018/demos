DOMAIN ?= apps.acme.com
REPO ?= github.com/scottd018/demos
GROUP ?= demo
VERSION ?= v1alpha1

define run_clean
	rm -rf $(1)/*; \
	cd $(1); \
	rm -rf .dockerignore .gitignore
endef

define run_init
	cd $(1); \
	$(2) init \
		--domain $(DOMAIN) \
		--project-name $(3) \
		$(4)
endef

define run_init_ob
	cd $(1); \
	$(2) init \
		--repo $(REPO) \
		--workload-config $(3) \
		--skip-go-version-check
endef

define run_create
	cd $(1); \
	$(2) create api \
	--group $(GROUP) \
	--version $(VERSION) \
	--kind $(3) \
	$(4)
endef

define run_create_ob
	cd $(1); \
	$(2) create api \
	--controller \
	--resource \
	--workload-config $(3)
endef

#
# KUBEBUILDER
#
KUBEBUILDER_DIR ?= kubebuilder
KUBEBUILDER_CMD ?= kubebuilder
KUBEBUILDER_KIND ?= Kubebuilder
KUBEBUILDER_PROJECT ?= kubebuilder-demo
kubebuilder-clean:
	$(call run_clean,$(KUBEBUILDER_DIR))

kubebuilder-init:
	$(call run_init,$(KUBEBUILDER_DIR),$(KUBEBUILDER_CMD),$(KUBEBUILDER_PROJECT))

kubebuilder-create:
	$(call run_create,$(KUBEBUILDER_DIR),$(KUBEBUILDER_CMD),$(KUBEBUILDER_KIND))

#
# OPERATOR SDK
#
OPERATOR_SDK_DIR ?= operator-sdk
OPERATOR_SDK_CMD ?= operator-sdk
OPERATOR_SDK_KIND ?= OperatorSDK
OPERATOR_SDK_PROJECT ?= operator-sdk-demo
operator-sdk-clean:
	$(call run_clean,$(OPERATOR_SDK_DIR))

operator-sdk-init:
	$(call run_init,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_PROJECT),--repo=$(REPO))

operator-sdk-create:
	$(call run_create,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_KIND),--controller --resource)

operator-sdk-init-ansible:
	$(call run_init,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_PROJECT),--plugins=ansible)

operator-sdk-create-ansible:
	$(call run_create,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_KIND),--generate-role)

operator-sdk-init-helm:
	$(call run_init,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_PROJECT),--plugins=helm)

operator-sdk-create-helm:
	$(call run_create,$(OPERATOR_SDK_DIR),$(OPERATOR_SDK_CMD),$(OPERATOR_SDK_KIND))

#
# OPERATOR BUILDER
#
OPERATOR_BUILDER_DIR ?= operator-builder
OPERATOR_BUILDER_CMD ?= operator-builder
OPERATOR_BUILDER_KIND ?= OperatorBuilder
OPERATOR_BUILDER_PROJECT ?= operator-builder-demo
OPERATOR_BUILDER_CONFIG ?= .workloadConfig/workload.yaml
operator-builder-clean:
	$(call run_clean,$(OPERATOR_BUILDER_DIR))

operator-builder-init:
	$(call run_init_ob,$(OPERATOR_BUILDER_DIR),$(OPERATOR_BUILDER_CMD),$(OPERATOR_BUILDER_CONFIG))

operator-builder-create:
	$(call run_create_ob,$(OPERATOR_BUILDER_DIR),$(OPERATOR_BUILDER_CMD),$(OPERATOR_BUILDER_CONFIG))
