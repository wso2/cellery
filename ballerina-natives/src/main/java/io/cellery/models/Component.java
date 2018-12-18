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
    List<API> apis;
    Set<Egress> egresses;
    String source;
    String service;
    int containerPort;
    int servicePort;

    public Component() {
        envVars = new HashMap<>();
        apis = new ArrayList<>();
        egresses = new HashSet<>();
        replicas = 1;
    }

    public void addApi(API api) {
        this.apis.add(api);
    }

    public void addEgress(Egress egress) {
        this.egresses.add(egress);
    }

    public void addEnv(String key, String value) {
        this.envVars.put(key, value);
    }
}
