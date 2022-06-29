resource "agile_logical_port" "example" {
  name            = "example"
  description     = "This Logical Port is created by terraform"
  tenant_id       = "11ade37a-79d0-482f-a7a0-6ad070e1d05d"
  fabric_id       = "f1429224-1860-4bdb-8cc8-98ccc0f5563a"
  logic_switch_id = "6c0a96d3-0789-47e6-9dbc-66ac5ba2e519"
  access_info {
    mode = "Uni"
    type = "Dot1q"
    vlan = 1218
    qinq {
      inner_vid_begin = 10
      inner_vid_end   = 10
      outer_vid_begin = 10
      outer_vid_end   = 10
      rewrite_action  = "PopDouble"
    }
    location {
      device_group_id = "e13784fb-499f-4c30-8f9c-e49e6c98fdbb"
      device_id       = "9e3a5bee-3d95-3bf7-90f5-09bd2177324b"
      port_id         = "589c87dd-7222-3c09-87b7-d09a236af285"
    }
    location {
      device_group_id = "e13784fb-499f-4c30-8f9c-e49e6c98fdbb"
      device_id       = "b4f6d9ed-0f1d-3f7a-82f1-a4a7ea4f84d4"
      port_id         = "4c142b5e-1858-33b2-a03e-71dcc3b37360"
    }
    subinterface_number = 18
  }
  additional {
    producer = "Terraform"
  }
}