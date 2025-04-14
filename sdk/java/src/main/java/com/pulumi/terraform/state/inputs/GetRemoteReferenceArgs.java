// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.terraform.state.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.terraform.state.inputs.WorkspacesArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class GetRemoteReferenceArgs extends com.pulumi.resources.InvokeArgs {

    public static final GetRemoteReferenceArgs Empty = new GetRemoteReferenceArgs();

    /**
     * The remote backend hostname to connect to.
     * 
     */
    @Import(name="hostname")
    private @Nullable Output<String> hostname;

    /**
     * @return The remote backend hostname to connect to.
     * 
     */
    public Optional<Output<String>> hostname() {
        return Optional.ofNullable(this.hostname);
    }

    /**
     * The name of the organization containing the targeted workspace(s).
     * 
     */
    @Import(name="organization", required=true)
    private Output<String> organization;

    /**
     * @return The name of the organization containing the targeted workspace(s).
     * 
     */
    public Output<String> organization() {
        return this.organization;
    }

    /**
     * The token used to authenticate with the remote backend.
     * 
     */
    @Import(name="token")
    private @Nullable Output<String> token;

    /**
     * @return The token used to authenticate with the remote backend.
     * 
     */
    public Optional<Output<String>> token() {
        return Optional.ofNullable(this.token);
    }

    @Import(name="workspaces", required=true)
    private Output<WorkspacesArgs> workspaces;

    public Output<WorkspacesArgs> workspaces() {
        return this.workspaces;
    }

    private GetRemoteReferenceArgs() {}

    private GetRemoteReferenceArgs(GetRemoteReferenceArgs $) {
        this.hostname = $.hostname;
        this.organization = $.organization;
        this.token = $.token;
        this.workspaces = $.workspaces;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(GetRemoteReferenceArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private GetRemoteReferenceArgs $;

        public Builder() {
            $ = new GetRemoteReferenceArgs();
        }

        public Builder(GetRemoteReferenceArgs defaults) {
            $ = new GetRemoteReferenceArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param hostname The remote backend hostname to connect to.
         * 
         * @return builder
         * 
         */
        public Builder hostname(@Nullable Output<String> hostname) {
            $.hostname = hostname;
            return this;
        }

        /**
         * @param hostname The remote backend hostname to connect to.
         * 
         * @return builder
         * 
         */
        public Builder hostname(String hostname) {
            return hostname(Output.of(hostname));
        }

        /**
         * @param organization The name of the organization containing the targeted workspace(s).
         * 
         * @return builder
         * 
         */
        public Builder organization(Output<String> organization) {
            $.organization = organization;
            return this;
        }

        /**
         * @param organization The name of the organization containing the targeted workspace(s).
         * 
         * @return builder
         * 
         */
        public Builder organization(String organization) {
            return organization(Output.of(organization));
        }

        /**
         * @param token The token used to authenticate with the remote backend.
         * 
         * @return builder
         * 
         */
        public Builder token(@Nullable Output<String> token) {
            $.token = token;
            return this;
        }

        /**
         * @param token The token used to authenticate with the remote backend.
         * 
         * @return builder
         * 
         */
        public Builder token(String token) {
            return token(Output.of(token));
        }

        public Builder workspaces(Output<WorkspacesArgs> workspaces) {
            $.workspaces = workspaces;
            return this;
        }

        public Builder workspaces(WorkspacesArgs workspaces) {
            return workspaces(Output.of(workspaces));
        }

        public GetRemoteReferenceArgs build() {
            $.hostname = Codegen.stringProp("hostname").output().arg($.hostname).def("app.terraform.io").getNullable();
            if ($.organization == null) {
                throw new MissingRequiredPropertyException("GetRemoteReferenceArgs", "organization");
            }
            if ($.workspaces == null) {
                throw new MissingRequiredPropertyException("GetRemoteReferenceArgs", "workspaces");
            }
            return $;
        }
    }

}
