// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.terraform.state.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class WorkspacesArgs extends com.pulumi.resources.ResourceArgs {

    public static final WorkspacesArgs Empty = new WorkspacesArgs();

    /**
     * The full name of one remote workspace. When configured, only the default workspace can be used. This option conflicts with prefix.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return The full name of one remote workspace. When configured, only the default workspace can be used. This option conflicts with prefix.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * A prefix used in the names of one or more remote workspaces, all of which can be used with this configuration. The full workspace names are used in HCP Terraform, and the short names (minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only the default workspace can be used. This option conflicts with name.
     * 
     */
    @Import(name="prefix")
    private @Nullable Output<String> prefix;

    /**
     * @return A prefix used in the names of one or more remote workspaces, all of which can be used with this configuration. The full workspace names are used in HCP Terraform, and the short names (minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only the default workspace can be used. This option conflicts with name.
     * 
     */
    public Optional<Output<String>> prefix() {
        return Optional.ofNullable(this.prefix);
    }

    private WorkspacesArgs() {}

    private WorkspacesArgs(WorkspacesArgs $) {
        this.name = $.name;
        this.prefix = $.prefix;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(WorkspacesArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private WorkspacesArgs $;

        public Builder() {
            $ = new WorkspacesArgs();
        }

        public Builder(WorkspacesArgs defaults) {
            $ = new WorkspacesArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name The full name of one remote workspace. When configured, only the default workspace can be used. This option conflicts with prefix.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name The full name of one remote workspace. When configured, only the default workspace can be used. This option conflicts with prefix.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param prefix A prefix used in the names of one or more remote workspaces, all of which can be used with this configuration. The full workspace names are used in HCP Terraform, and the short names (minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only the default workspace can be used. This option conflicts with name.
         * 
         * @return builder
         * 
         */
        public Builder prefix(@Nullable Output<String> prefix) {
            $.prefix = prefix;
            return this;
        }

        /**
         * @param prefix A prefix used in the names of one or more remote workspaces, all of which can be used with this configuration. The full workspace names are used in HCP Terraform, and the short names (minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only the default workspace can be used. This option conflicts with name.
         * 
         * @return builder
         * 
         */
        public Builder prefix(String prefix) {
            return prefix(Output.of(prefix));
        }

        public WorkspacesArgs build() {
            return $;
        }
    }

}
