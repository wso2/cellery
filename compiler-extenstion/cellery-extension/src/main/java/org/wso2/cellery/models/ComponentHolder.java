package org.wso2.cellery.models;

import java.util.ArrayList;
import java.util.List;

/**
 * Singleton Component Holder Class.
 */
public class ComponentHolder {


    private List<Component> components;

    private ComponentHolder() {
        components = new ArrayList<>();
    }

    public static ComponentHolder getInstance() {
        return SingletonHelper.INSTANCE;
    }

    public List<Component> getComponents() {
        return components;
    }

    public void addComponent(Component component) {
        this.components.add(component);
    }

    private static class SingletonHelper {
        private static final ComponentHolder INSTANCE = new ComponentHolder();
    }
}
