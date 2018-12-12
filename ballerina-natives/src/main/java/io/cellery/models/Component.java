package io.cellery.models;

import lombok.Data;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Cell component.
 */
@Data
public class Component {
    String name;
    int replicas;
    Map<String, String> envVars;
    List<API> apis;
    String source;
    String service;
    int containerPort;
    int servicePort;

    public Component() {
        envVars = new HashMap<>();
        apis = new ArrayList<>();
    }

    public void addApi(API api) {
        this.apis.add(api);
    }

    public void addEnv(String key, String value) {
        this.envVars.put(key, value);
    }
}
