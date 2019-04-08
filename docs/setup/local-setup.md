### 1. Local setup
This will setup the local environment, by creating a virtual machine with pre-installed kubeadm and cellery runtime. 

1) As mentioned above this can be installed with [interactively](../../README.md#interactive-mode-setup) by selecting Create > Local > Basic or Complete options 
or executing [inline command](../../README.md#inline-command-mode-setup) with `cellery setup create local [--complete]`.

2) Add below to /etc/host entries to access cellery hosts.
    ```
    192.168.56.10 wso2-apim cellery-dashboard wso2sp-observability-api wso2-apim-gateway cellery-k8s-metrics idp.cellery-system pet-store.com hello-world.com my-hello-world.com
    ```
3) As the installation process is completed, you can [quick start with cellery](../../README.md#quick-start-with-cellery).
