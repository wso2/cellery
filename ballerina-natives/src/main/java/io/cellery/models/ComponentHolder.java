package io.cellery.models;

import org.ballerinalang.util.exceptions.BallerinaException;

import java.util.HashMap;
import java.util.Map;

/**
 * Singleton Component Holder Class.
 */
public class ComponentHolder {

    private Map<String, Component> componentNameToComponentMap;

    private ComponentHolder() {
        componentNameToComponentMap = new HashMap<>();
    }

    public static ComponentHolder getInstance() {
        return SingletonHelper.INSTANCE;
    }

    public Map<String, Component> getComponentNameToComponentMap() {
        return componentNameToComponentMap;
    }

    public void addComponent(Component component) {
        if (componentNameToComponentMap.get(component.getName()) != null) {
            throw new BallerinaException("Two components with same name exists " + component.getName());
        }
        this.componentNameToComponentMap.put(component.getName(), component);

    }

    public void addAPI(String componentName, API api) {
        Component temp = componentNameToComponentMap.remove(componentName);
        if (temp == null) {
            throw new BallerinaException("Invalid component name " + componentName);
        }
        api.setBackend(temp.getService());
        temp.addApi(api);
        componentNameToComponentMap.put(componentName, temp);
    }

    private static class SingletonHelper {
        private static final ComponentHolder INSTANCE = new ComponentHolder();
    }
}
