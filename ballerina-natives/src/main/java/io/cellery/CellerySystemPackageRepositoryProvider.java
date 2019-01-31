package io.cellery;

import org.ballerinalang.annotation.JavaSPIService;
import org.ballerinalang.spi.SystemPackageRepositoryProvider;
import org.wso2.ballerinalang.compiler.packaging.repo.JarRepo;
import org.wso2.ballerinalang.compiler.packaging.repo.Repo;

/**
 * This represents the Ballerina regex package repository provider.
 */
@JavaSPIService("org.ballerinalang.spi.SystemPackageRepositoryProvider")
public class CellerySystemPackageRepositoryProvider implements SystemPackageRepositoryProvider {
    @Override
    public Repo loadRepository() {
        return new JarRepo(SystemPackageRepositoryProvider.getClassUri(this));
    }
}
