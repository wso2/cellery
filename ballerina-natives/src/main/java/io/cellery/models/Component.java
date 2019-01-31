package io.cellery.models;

import lombok.Data;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Cell component.
 */
@Data
public class Component {
    String name;
    int replicas;
    Map<String, String> envVars;
    Map<String, String> labels;
    List<API> apis;
    Set<Egress> egresses;
    String source;
    String service;
    Boolean isStub;
    Map<Integer, Integer> containerPortToServicePortMap;

    public Component() {
        envVars = new HashMap<>();
        labels = new HashMap<>();
        apis = new ArrayList<>();
        egresses = new HashSet<>();
        containerPortToServicePortMap = new HashMap<>();
        replicas = 1;
        isStub = false;
    }

    void addApi(API api) {
        this.apis.add(api);
    }

    public void addPorts(int containerMap, int servicePort) {
        this.containerPortToServicePortMap.put(containerMap, servicePort);
    }

    public void addEgress(Egress egress) {
        this.egresses.add(egress);
    }

    public void addEnv(String key, String value) {
        this.envVars.put(key, value);
    }

    public void addLabel(String key, String value) {
        this.labels.put(key, value);
    }
}
