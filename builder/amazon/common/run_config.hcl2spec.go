// Code generated by "hcl2-schema"; DO NOT EDIT.\n

package common

import (
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

func (*RunConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"AssociatePublicIpAddress":          &hcldec.AttrSpec{Name: "associate_public_ip_address", Type: cty.Bool, Required: false},
		"AvailabilityZone":                  &hcldec.AttrSpec{Name: "availability_zone", Type: cty.String, Required: false},
		"BlockDurationMinutes":              &hcldec.AttrSpec{Name: "block_duration_minutes", Type: cty.Number, Required: false},
		"DisableStopInstance":               &hcldec.AttrSpec{Name: "disable_stop_instance", Type: cty.Bool, Required: false},
		"EbsOptimized":                      &hcldec.AttrSpec{Name: "ebs_optimized", Type: cty.Bool, Required: false},
		"EnableT2Unlimited":                 &hcldec.AttrSpec{Name: "enable_t2_unlimited", Type: cty.Bool, Required: false},
		"IamInstanceProfile":                &hcldec.AttrSpec{Name: "iam_instance_profile", Type: cty.String, Required: false},
		"InstanceInitiatedShutdownBehavior": &hcldec.AttrSpec{Name: "shutdown_behavior", Type: cty.String, Required: false},
		"InstanceType":                      &hcldec.AttrSpec{Name: "instance_type", Type: cty.String, Required: false},
		"SecurityGroupId":                   &hcldec.AttrSpec{Name: "security_group_id", Type: cty.String, Required: false},
		"SecurityGroupIds":                  &hcldec.AttrSpec{Name: "security_group_ids", Type: cty.List(cty.String), Required: false},
		"SourceAmi":                         &hcldec.AttrSpec{Name: "source_ami", Type: cty.String, Required: false},
		"SpotInstanceTypes":                 &hcldec.AttrSpec{Name: "spot_instance_types", Type: cty.List(cty.String), Required: false},
		"SpotPrice":                         &hcldec.AttrSpec{Name: "spot_price", Type: cty.String, Required: false},
		"SpotPriceAutoProduct":              &hcldec.AttrSpec{Name: "spot_price_auto_product", Type: cty.String, Required: false},
		"SubnetId":                          &hcldec.AttrSpec{Name: "subnet_id", Type: cty.String, Required: false},
		"TemporaryKeyPairName":              &hcldec.AttrSpec{Name: "temporary_key_pair_name", Type: cty.String, Required: false},
		"TemporarySGSourceCidrs":            &hcldec.AttrSpec{Name: "temporary_security_group_source_cidrs", Type: cty.List(cty.String), Required: false},
		"UserData":                          &hcldec.AttrSpec{Name: "user_data", Type: cty.String, Required: false},
		"UserDataFile":                      &hcldec.AttrSpec{Name: "user_data_file", Type: cty.String, Required: false},
		"VpcId":                             &hcldec.AttrSpec{Name: "vpc_id", Type: cty.String, Required: false},
		"WindowsPasswordTimeout":            &hcldec.AttrSpec{Name: "windows_password_timeout", Type: cty.String, Required: false},
		"SSHInterface":                      &hcldec.AttrSpec{Name: "ssh_interface", Type: cty.String, Required: false},
		"security_group_filter":             &hcldec.BlockObjectSpec{TypeName: "SecurityGroupFilterOptions", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&RunConfig{}).SecurityGroupFilter.HCL2Spec())},
		"source_ami_filter":                 &hcldec.BlockObjectSpec{TypeName: "AmiFilterOptions", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&RunConfig{}).SourceAmiFilter.HCL2Spec())},
		"subnet_filter":                     &hcldec.BlockObjectSpec{TypeName: "SubnetFilterOptions", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&RunConfig{}).SubnetFilter.HCL2Spec())},
		"vpc_filter":                        &hcldec.BlockObjectSpec{TypeName: "VpcFilterOptions", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&RunConfig{}).VpcFilter.HCL2Spec())},
	}
	return s
}

func (*AmiFilterOptions) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"MostRecent": &hcldec.AttrSpec{Name: "most_recent", Type: cty.Bool, Required: false},
		"owners":     &hcldec.BlockObjectSpec{TypeName: "[]*string", LabelNames: []string(nil), Nested: hcldec.ObjectSpec((&AmiFilterOptions{}).Owners.HCL2Spec())},
	}
	return s
}

func (*SecurityGroupFilterOptions) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{}
	return s
}

func (*SubnetFilterOptions) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"MostFree": &hcldec.AttrSpec{Name: "most_free", Type: cty.Bool, Required: false},
		"Random":   &hcldec.AttrSpec{Name: "random", Type: cty.Bool, Required: false},
	}
	return s
}

func (*VpcFilterOptions) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{}
	return s
}
