/*
Copyright (c) 2022-2023 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var serializationTests = []struct {
	name   string
	file   string
	data   interface{}
	goType reflect.Type
}{
	{
		name:   "vminfo",
		file:   "./testdata/vminfo.json",
		data:   &vmInfoObjForTests,
		goType: reflect.TypeOf(VirtualMachineConfigInfo{}),
	},
	{
		name:   "retrieveResult",
		file:   "./testdata/retrieveResult.json",
		data:   &retrieveResultForTests,
		goType: reflect.TypeOf(RetrieveResult{}),
	},
}

func TestSerialization(t *testing.T) {
	for _, test := range serializationTests {
		t.Run(test.name+" Decode", func(t *testing.T) {
			f, err := os.Open(test.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			dec := NewGovmomiDecoder(f)

			data := reflect.New(test.goType).Interface()
			if err := dec.Decode(data); err != nil {
				t.Fatal(err)
			}

			a, e := data, test.data

			if diff := cmp.Diff(a, e); diff != "" {
				t.Errorf("mismatched %v: %s", test.name, diff)
			}
		})

		t.Run(test.name+" Encode", func(t *testing.T) {
			expJSON, err := os.ReadFile(test.file)
			if err != nil {
				t.Fatal(err)
			}

			var w bytes.Buffer
			_ = w
			enc := NewGovmomiEncoder(&w)

			if err := enc.Encode(test.data); err != nil {
				t.Fatal(err)
			}

			expected, actual := string(expJSON), w.String()
			assert.JSONEq(t, expected, actual)
		})
	}
}

var vmInfoObjForTests = VirtualMachineConfigInfo{
	ChangeVersion:         "2022-12-12T11:48:35.473645Z",
	Modified:              mustParseTime(time.RFC3339, "1970-01-01T00:00:00Z"),
	Name:                  "test",
	GuestFullName:         "VMware Photon OS (64-bit)",
	Version:               "vmx-20",
	Uuid:                  "422ca90b-853b-1101-3350-759f747730cc",
	CreateDate:            addrOfMustParseTime(time.RFC3339, "2022-12-12T11:47:24.685785Z"),
	InstanceUuid:          "502cc2a5-1f06-2890-6d70-ba2c55c5c2b7",
	NpivTemporaryDisabled: addrOfBool(true),
	LocationId:            "Earth",
	Template:              false,
	GuestId:               "vmwarePhoton64Guest",
	AlternateGuestName:    "",
	Annotation:            "Hello, world.",
	Files: VirtualMachineFileInfo{
		VmPathName:        "[datastore1] test/test.vmx",
		SnapshotDirectory: "[datastore1] test/",
		SuspendDirectory:  "[datastore1] test/",
		LogDirectory:      "[datastore1] test/",
	},
	Tools: &ToolsConfigInfo{
		ToolsVersion:            1,
		AfterPowerOn:            addrOfBool(true),
		AfterResume:             addrOfBool(true),
		BeforeGuestStandby:      addrOfBool(true),
		BeforeGuestShutdown:     addrOfBool(true),
		BeforeGuestReboot:       nil,
		ToolsUpgradePolicy:      "manual",
		SyncTimeWithHostAllowed: addrOfBool(true),
		SyncTimeWithHost:        addrOfBool(false),
		LastInstallInfo: &ToolsConfigInfoToolsLastInstallInfo{
			Counter: 0,
		},
	},
	Flags: VirtualMachineFlagInfo{
		EnableLogging:            addrOfBool(true),
		UseToe:                   addrOfBool(false),
		RunWithDebugInfo:         addrOfBool(false),
		MonitorType:              "release",
		HtSharing:                "any",
		SnapshotDisabled:         addrOfBool(false),
		SnapshotLocked:           addrOfBool(false),
		DiskUuidEnabled:          addrOfBool(false),
		SnapshotPowerOffBehavior: "powerOff",
		RecordReplayEnabled:      addrOfBool(false),
		FaultToleranceType:       "unset",
		CbrcCacheEnabled:         addrOfBool(false),
		VvtdEnabled:              addrOfBool(false),
		VbsEnabled:               addrOfBool(false),
	},
	DefaultPowerOps: VirtualMachineDefaultPowerOpInfo{
		PowerOffType:        "soft",
		SuspendType:         "hard",
		ResetType:           "soft",
		DefaultPowerOffType: "soft",
		DefaultSuspendType:  "hard",
		DefaultResetType:    "soft",
		StandbyAction:       "checkpoint",
	},
	RebootPowerOff: addrOfBool(false),
	Hardware: VirtualHardware{
		NumCPU:              1,
		NumCoresPerSocket:   1,
		AutoCoresPerSocket:  addrOfBool(true),
		MemoryMB:            2048,
		VirtualICH7MPresent: addrOfBool(false),
		VirtualSMCPresent:   addrOfBool(false),
		Device: []BaseVirtualDevice{
			&VirtualIDEController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 200,
						DeviceInfo: &Description{
							Label:   "IDE 0",
							Summary: "IDE 0",
						},
					},
					BusNumber: 0,
				},
			},
			&VirtualIDEController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 201,
						DeviceInfo: &Description{
							Label:   "IDE 1",
							Summary: "IDE 1",
						},
					},
					BusNumber: 1,
				},
			},
			&VirtualPS2Controller{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 300,
						DeviceInfo: &Description{
							Label:   "PS2 controller 0",
							Summary: "PS2 controller 0",
						},
					},
					BusNumber: 0,
					Device:    []int32{600, 700},
				},
			},
			&VirtualPCIController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 100,
						DeviceInfo: &Description{
							Label:   "PCI controller 0",
							Summary: "PCI controller 0",
						},
					},
					BusNumber: 0,
					Device:    []int32{500, 12000, 14000, 1000, 15000, 4000},
				},
			},
			&VirtualSIOController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 400,
						DeviceInfo: &Description{
							Label:   "SIO controller 0",
							Summary: "SIO controller 0",
						},
					},
					BusNumber: 0,
				},
			},
			&VirtualKeyboard{
				VirtualDevice: VirtualDevice{
					Key: 600,
					DeviceInfo: &Description{
						Label:   "Keyboard",
						Summary: "Keyboard",
					},
					ControllerKey: 300,
					UnitNumber:    addrOfInt32(0),
				},
			},
			&VirtualPointingDevice{
				VirtualDevice: VirtualDevice{
					Key:        700,
					DeviceInfo: &Description{Label: "Pointing device", Summary: "Pointing device; Device"},
					Backing: &VirtualPointingDeviceDeviceBackingInfo{
						VirtualDeviceDeviceBackingInfo: VirtualDeviceDeviceBackingInfo{
							UseAutoDetect: addrOfBool(false),
						},
						HostPointingDevice: "autodetect",
					},
					ControllerKey: 300,
					UnitNumber:    addrOfInt32(1),
				},
			},
			&VirtualMachineVideoCard{
				VirtualDevice: VirtualDevice{
					Key:           500,
					DeviceInfo:    &Description{Label: "Video card ", Summary: "Video card"},
					ControllerKey: 100,
					UnitNumber:    addrOfInt32(0),
				},
				VideoRamSizeInKB:       4096,
				NumDisplays:            1,
				UseAutoDetect:          addrOfBool(false),
				Enable3DSupport:        addrOfBool(false),
				Use3dRenderer:          "automatic",
				GraphicsMemorySizeInKB: 262144,
			},
			&VirtualMachineVMCIDevice{
				VirtualDevice: VirtualDevice{
					Key: 12000,
					DeviceInfo: &Description{
						Label: "VMCI device",
						Summary: "Device on the virtual machine PCI " +
							"bus that provides support for the " +
							"virtual machine communication interface",
					},
					ControllerKey: 100,
					UnitNumber:    addrOfInt32(17),
				},
				Id:                             -1,
				AllowUnrestrictedCommunication: addrOfBool(false),
				FilterEnable:                   addrOfBool(true),
			},
			&ParaVirtualSCSIController{
				VirtualSCSIController: VirtualSCSIController{
					VirtualController: VirtualController{
						VirtualDevice: VirtualDevice{
							Key: 1000,
							DeviceInfo: &Description{
								Label:   "SCSI controller 0",
								Summary: "VMware paravirtual SCSI",
							},
							ControllerKey: 100,
							UnitNumber:    addrOfInt32(3),
						},
						Device: []int32{2000},
					},
					HotAddRemove:       addrOfBool(true),
					SharedBus:          "noSharing",
					ScsiCtlrUnitNumber: 7,
				},
			},
			&VirtualAHCIController{
				VirtualSATAController: VirtualSATAController{
					VirtualController: VirtualController{
						VirtualDevice: VirtualDevice{
							Key: 15000,
							DeviceInfo: &Description{
								Label:   "SATA controller 0",
								Summary: "AHCI",
							},
							ControllerKey: 100,
							UnitNumber:    addrOfInt32(24),
						},
						Device: []int32{16000},
					},
				},
			},
			&VirtualCdrom{
				VirtualDevice: VirtualDevice{
					Key: 16000,
					DeviceInfo: &Description{
						Label:   "CD/DVD drive 1",
						Summary: "Remote device",
					},
					Backing: &VirtualCdromRemotePassthroughBackingInfo{
						VirtualDeviceRemoteDeviceBackingInfo: VirtualDeviceRemoteDeviceBackingInfo{
							UseAutoDetect: addrOfBool(false),
						},
					},
					Connectable:   &VirtualDeviceConnectInfo{AllowGuestControl: true, Status: "untried"},
					ControllerKey: 15000,
					UnitNumber:    addrOfInt32(0),
				},
			},
			&VirtualDisk{
				VirtualDevice: VirtualDevice{
					Key: 2000,
					DeviceInfo: &Description{
						Label:   "Hard disk 1",
						Summary: "4,194,304 KB",
					},
					Backing: &VirtualDiskFlatVer2BackingInfo{
						VirtualDeviceFileBackingInfo: VirtualDeviceFileBackingInfo{
							BackingObjectId: "1",
							FileName:        "[datastore1] test/test.vmdk",
							Datastore: &ManagedObjectReference{
								Type:  "Datastore",
								Value: "datastore-21",
							},
						},
						DiskMode:               "persistent",
						Split:                  addrOfBool(false),
						WriteThrough:           addrOfBool(false),
						ThinProvisioned:        addrOfBool(false),
						EagerlyScrub:           addrOfBool(false),
						Uuid:                   "6000C298-df15-fe89-ddcb-8ea33329595d",
						ContentId:              "e4e1a794c6307ce7906a3973fffffffe",
						ChangeId:               "",
						Parent:                 nil,
						DeltaDiskFormat:        "",
						DigestEnabled:          addrOfBool(false),
						DeltaGrainSize:         0,
						DeltaDiskFormatVariant: "",
						Sharing:                "sharingNone",
						KeyId:                  nil,
					},
					ControllerKey: 1000,
					UnitNumber:    addrOfInt32(0),
				},
				CapacityInKB:    4194304,
				CapacityInBytes: 4294967296,
				Shares:          &SharesInfo{Shares: 1000, Level: "normal"},
				StorageIOAllocation: &StorageIOAllocationInfo{
					Limit:       addrOfInt64(-1),
					Shares:      &SharesInfo{Shares: 1000, Level: "normal"},
					Reservation: addrOfInt32(0),
				},
				DiskObjectId:               "1-2000",
				NativeUnmanagedLinkedClone: addrOfBool(false),
			},
			&VirtualVmxnet3{
				VirtualVmxnet: VirtualVmxnet{
					VirtualEthernetCard: VirtualEthernetCard{
						VirtualDevice: VirtualDevice{
							Key: 4000,
							DeviceInfo: &Description{
								Label:   "Network adapter 1",
								Summary: "VM Network",
							},
							Backing: &VirtualEthernetCardNetworkBackingInfo{
								VirtualDeviceDeviceBackingInfo: VirtualDeviceDeviceBackingInfo{
									DeviceName:    "VM Network",
									UseAutoDetect: addrOfBool(false),
								},
								Network: &ManagedObjectReference{
									Type:  "Network",
									Value: "network-27",
								},
							},
							Connectable: &VirtualDeviceConnectInfo{
								MigrateConnect: "unset",
								StartConnected: true,
								Status:         "untried",
							},
							ControllerKey: 100,
							UnitNumber:    addrOfInt32(7),
						},
						AddressType:      "assigned",
						MacAddress:       "00:50:56:ac:4d:ed",
						WakeOnLanEnabled: addrOfBool(true),
						ResourceAllocation: &VirtualEthernetCardResourceAllocation{
							Reservation: addrOfInt64(0),
							Share: SharesInfo{
								Shares: 50,
								Level:  "normal",
							},
							Limit: addrOfInt64(-1),
						},
						UptCompatibilityEnabled: addrOfBool(true),
					},
				},
				Uptv2Enabled: addrOfBool(false),
			},
			&VirtualUSBXHCIController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						Key: 14000,
						DeviceInfo: &Description{
							Label:   "USB xHCI controller ",
							Summary: "USB xHCI controller",
						},
						SlotInfo: &VirtualDevicePciBusSlotInfo{
							PciSlotNumber: -1,
						},
						ControllerKey: 100,
						UnitNumber:    addrOfInt32(23),
					},
				},

				AutoConnectDevices: addrOfBool(false),
			},
		},
		MotherboardLayout:   "i440bxHostBridge",
		SimultaneousThreads: 1,
	},
	CpuAllocation: &ResourceAllocationInfo{
		Reservation:           addrOfInt64(0),
		ExpandableReservation: addrOfBool(false),
		Limit:                 addrOfInt64(-1),
		Shares: &SharesInfo{
			Shares: 1000,
			Level:  SharesLevelNormal,
		},
	},
	MemoryAllocation: &ResourceAllocationInfo{
		Reservation:           addrOfInt64(0),
		ExpandableReservation: addrOfBool(false),
		Limit:                 addrOfInt64(-1),
		Shares: &SharesInfo{
			Shares: 20480,
			Level:  SharesLevelNormal,
		},
	},
	LatencySensitivity: &LatencySensitivity{
		Level: LatencySensitivitySensitivityLevelNormal,
	},
	MemoryHotAddEnabled: addrOfBool(false),
	CpuHotAddEnabled:    addrOfBool(false),
	CpuHotRemoveEnabled: addrOfBool(false),
	ExtraConfig: []BaseOptionValue{
		&OptionValue{Key: "nvram", Value: "test.nvram"},
		&OptionValue{Key: "svga.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge0.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge4.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge4.virtualDev", Value: "pcieRootPort"},
		&OptionValue{Key: "pciBridge4.functions", Value: "8"},
		&OptionValue{Key: "pciBridge5.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge5.virtualDev", Value: "pcieRootPort"},
		&OptionValue{Key: "pciBridge5.functions", Value: "8"},
		&OptionValue{Key: "pciBridge6.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge6.virtualDev", Value: "pcieRootPort"},
		&OptionValue{Key: "pciBridge6.functions", Value: "8"},
		&OptionValue{Key: "pciBridge7.present", Value: "TRUE"},
		&OptionValue{Key: "pciBridge7.virtualDev", Value: "pcieRootPort"},
		&OptionValue{Key: "pciBridge7.functions", Value: "8"},
		&OptionValue{Key: "hpet0.present", Value: "TRUE"},
		&OptionValue{Key: "RemoteDisplay.maxConnections", Value: "-1"},
		&OptionValue{Key: "sched.cpu.latencySensitivity", Value: "normal"},
		&OptionValue{Key: "vmware.tools.internalversion", Value: "0"},
		&OptionValue{Key: "vmware.tools.requiredversion", Value: "12352"},
		&OptionValue{Key: "migrate.hostLogState", Value: "none"},
		&OptionValue{Key: "migrate.migrationId", Value: "0"},
		&OptionValue{Key: "migrate.hostLog", Value: "test-36f94569.hlog"},
		&OptionValue{
			Key:   "viv.moid",
			Value: "c5b34aa9-d962-4a74-b7d2-b83ec683ba1b:vm-28:lIgQ2t7v24n2nl3N7K3m6IHW2OoPF4CFrJd5N+Tdfio=",
		},
	},
	DatastoreUrl: []VirtualMachineConfigInfoDatastoreUrlPair{
		{
			Name: "datastore1",
			Url:  "/vmfs/volumes/63970ed8-4abddd2a-62d7-02003f49c37d",
		},
	},
	SwapPlacement: "inherit",
	BootOptions: &VirtualMachineBootOptions{
		EnterBIOSSetup:       addrOfBool(false),
		EfiSecureBootEnabled: addrOfBool(false),
		BootDelay:            1,
		BootRetryEnabled:     addrOfBool(false),
		BootRetryDelay:       10000,
		NetworkBootProtocol:  "ipv4",
	},
	FtInfo:                       nil,
	RepConfig:                    nil,
	VAppConfig:                   nil,
	VAssertsEnabled:              addrOfBool(false),
	ChangeTrackingEnabled:        addrOfBool(false),
	Firmware:                     "bios",
	MaxMksConnections:            -1,
	GuestAutoLockEnabled:         addrOfBool(true),
	ManagedBy:                    nil,
	MemoryReservationLockedToMax: addrOfBool(false),
	InitialOverhead: &VirtualMachineConfigInfoOverheadInfo{
		InitialMemoryReservation: 214446080,
		InitialSwapReservation:   2541883392,
	},
	NestedHVEnabled: addrOfBool(false),
	VPMCEnabled:     addrOfBool(false),
	ScheduledHardwareUpgradeInfo: &ScheduledHardwareUpgradeInfo{
		UpgradePolicy:                  "never",
		ScheduledHardwareUpgradeStatus: "none",
	},
	ForkConfigInfo:         nil,
	VFlashCacheReservation: 0,
	VmxConfigChecksum: []uint8{
		0x69, 0xf7, 0xa7, 0x9e,
		0xd1, 0xc2, 0x21, 0x4b,
		0x6c, 0x20, 0x77, 0x0a,
		0x94, 0x94, 0x99, 0xee,
		0x17, 0x5d, 0xdd, 0xa3,
	},
	MessageBusTunnelEnabled: addrOfBool(false),
	GuestIntegrityInfo: &VirtualMachineGuestIntegrityInfo{
		Enabled: addrOfBool(false),
	},
	MigrateEncryption: "opportunistic",
	SgxInfo: &VirtualMachineSgxInfo{
		FlcMode:            "unlocked",
		RequireAttestation: addrOfBool(false),
	},
	ContentLibItemInfo:      nil,
	FtEncryptionMode:        "ftEncryptionOpportunistic",
	GuestMonitoringModeInfo: &VirtualMachineGuestMonitoringModeInfo{},
	SevEnabled:              addrOfBool(false),
	NumaInfo: &VirtualMachineVirtualNumaInfo{
		AutoCoresPerNumaNode:    addrOfBool(true),
		VnumaOnCpuHotaddExposed: addrOfBool(false),
	},
	PmemFailoverEnabled:          addrOfBool(false),
	VmxStatsCollectionEnabled:    addrOfBool(true),
	VmOpNotificationToAppEnabled: addrOfBool(false),
	VmOpNotificationTimeout:      -1,
	DeviceSwap: &VirtualMachineVirtualDeviceSwap{
		LsiToPvscsi: &VirtualMachineVirtualDeviceSwapDeviceSwapInfo{
			Enabled:    addrOfBool(true),
			Applicable: addrOfBool(false),
			Status:     "none",
		},
	},
	Pmem:         nil,
	DeviceGroups: &VirtualMachineVirtualDeviceGroups{},
}

var retrieveResultForTests = RetrieveResult{
	Token: "",
	Objects: []ObjectContent{

		{

			DynamicData: DynamicData{},
			Obj: ManagedObjectReference{

				Type:  "Folder",
				Value: "group-d1",
			},
			PropSet: []DynamicProperty{
				{

					Name: "alarmActionsEnabled",
					Val:  true,
				},
				{

					Name: "availableField",
					Val: ArrayOfCustomFieldDef{

						CustomFieldDef: []CustomFieldDef{},
					},
				},

				{

					Name: "childEntity",
					Val: ArrayOfManagedObjectReference{
						ManagedObjectReference: []ManagedObjectReference{},
					},
				},
				{
					Name: "childType",
					Val: ArrayOfString{
						String: []string{
							"Folder",
							"Datacenter"},
					},
				},
				{
					Name: "configIssue",
					Val: ArrayOfEvent{
						Event: []BaseEvent{},
					},
				},
				{
					Name: "configStatus",
					Val:  ManagedEntityStatusGray},
				{
					Name: "customValue",
					Val: ArrayOfCustomFieldValue{
						CustomFieldValue: []BaseCustomFieldValue{},
					},
				},
				{
					Name: "declaredAlarmState",
					Val: ArrayOfAlarmState{
						AlarmState: []AlarmState{
							{
								Key: "alarm-328.group-d1",
								Entity: ManagedObjectReference{
									Type:  "Folder",
									Value: "group-d1"},
								Alarm: ManagedObjectReference{
									Type:  "Alarm",
									Value: "alarm-328"},
								OverallStatus: "gray",
								Time:          time.Date(2023, time.January, 14, 8, 57, 35, 279575000, time.UTC),
								Acknowledged:  addrOfBool(false),
							},
							{
								Key: "alarm-327.group-d1",
								Entity: ManagedObjectReference{
									Type:  "Folder",
									Value: "group-d1"},
								Alarm: ManagedObjectReference{
									Type:  "Alarm",
									Value: "alarm-327"},
								OverallStatus: "green",
								Time:          time.Date(2023, time.January, 14, 8, 56, 40, 83607000, time.UTC),
								Acknowledged:  addrOfBool(false),
								EventKey:      756,
							},
							{
								DynamicData: DynamicData{},
								Key:         "alarm-326.group-d1",
								Entity: ManagedObjectReference{
									Type:  "Folder",
									Value: "group-d1"},
								Alarm: ManagedObjectReference{
									Type:  "Alarm",
									Value: "alarm-326"},
								OverallStatus: "green",
								Time: time.Date(2023,
									time.January,
									14,
									8,
									56,
									35,
									82616000,
									time.UTC),
								Acknowledged: addrOfBool(false),
								EventKey:     751,
							},
						},
					},
				},
				{
					Name: "disabledMethod",
					Val: ArrayOfString{
						String: []string{},
					},
				},
				{
					Name: "effectiveRole",
					Val: ArrayOfInt{
						Int: []int32{-1},
					},
				},
				{
					Name: "name",
					Val:  "Datacenters"},
				{
					Name: "overallStatus",
					Val:  ManagedEntityStatusGray},
				{
					Name: "permission",
					Val: ArrayOfPermission{
						Permission: []Permission{
							{
								Entity: &ManagedObjectReference{
									Value: "group-d1",
									Type:  "Folder",
								},
								Principal: "VSPHERE.LOCAL\\vmware-vsm-2bd917c6-e084-4d1f-988d-a68f7525cc94",
								Group:     false,
								RoleId:    1034,
								Propagate: true},
							{
								Entity: &ManagedObjectReference{
									Value: "group-d1",
									Type:  "Folder",
								},
								Principal: "VSPHERE.LOCAL\\topologysvc-2bd917c6-e084-4d1f-988d-a68f7525cc94",
								Group:     false,
								RoleId:    1024,
								Propagate: true},
							{
								Entity: &ManagedObjectReference{
									Value: "group-d1",
									Type:  "Folder",
								},
								Principal: "VSPHERE.LOCAL\\vpxd-extension-2bd917c6-e084-4d1f-988d-a68f7525cc94",
								Group:     false,
								RoleId:    -1,
								Propagate: true},
						},
					},
				},
				{
					Name: "recentTask",
					Val: ArrayOfManagedObjectReference{
						ManagedObjectReference: []ManagedObjectReference{
							{
								Type:  "Task",
								Value: "task-186"},
							{
								Type:  "Task",
								Value: "task-187"},
							{
								Type:  "Task",
								Value: "task-188"},
						},
					},
				},
				{
					Name: "tag",
					Val: ArrayOfTag{
						Tag: []Tag{},
					},
				},
				{
					Name: "triggeredAlarmState",
					Val: ArrayOfAlarmState{
						AlarmState: []AlarmState{},
					},
				},
				{
					Name: "value",
					Val: ArrayOfCustomFieldValue{
						CustomFieldValue: []BaseCustomFieldValue{},
					},
				},
			},
			MissingSet: nil,
		},
	},
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func addrOfMustParseTime(layout, value string) *time.Time {
	t := mustParseTime(layout, value)
	return &t
}

func addrOfBool(v bool) *bool {
	return &v
}

func addrOfInt32(v int32) *int32 {
	return &v
}

func addrOfInt64(v int64) *int64 {
	return &v
}
