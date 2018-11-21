package org.wso2.cellery;

import org.wso2.ballerinalang.compiler.tree.expressions.BLangRecordLiteral;
import org.wso2.ballerinalang.compiler.tree.expressions.BLangSimpleVarRef;
import org.wso2.cellery.models.API;

import java.util.HashMap;
import java.util.Map;

/**
 * Parse Cell Components.
 */
class CelleryParser {

    /**
     * Create {@link API} object from key value pairs.
     *
     * @param apiKeyValues Record values
     * @return {@link API} object
     */
    static API parseAPI(BLangRecordLiteral.BLangRecordKeyValue apiKeyValues) {
        API api = new API();
        ((BLangRecordLiteral) apiKeyValues.valueExpr).keyValuePairs.forEach(bLangRecordKeyValue -> {
            switch (APIEnum.valueOf(((BLangSimpleVarRef) bLangRecordKeyValue.key.expr).variableName.value)) {
                case context:
                    api.setContext(bLangRecordKeyValue.valueExpr.toString());
                    break;
                case global:
                    api.setGlobal(Boolean.parseBoolean(bLangRecordKeyValue.valueExpr.toString()));
                    break;
                default:
                    break;
            }
        });
        return api;
    }

    /**
     * Create Map object from key value pairs.
     *
     * @param envKeyValues Record values
     * @return Map object
     */
    static Map<String, String> parseEnv(BLangRecordLiteral.BLangRecordKeyValue envKeyValues) {
        Map<String, String> envVars = new HashMap<>();
        ((BLangRecordLiteral) envKeyValues.valueExpr).keyValuePairs.forEach(bLangRecordKeyValue -> {
            envVars.put(bLangRecordKeyValue.key.expr.toString(), bLangRecordKeyValue.valueExpr.toString());
        });
        return envVars;
    }

    private enum APIEnum {
        context,
        global,
        definitions,
    }
}
