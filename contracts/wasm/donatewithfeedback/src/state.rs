// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

#![allow(dead_code)]
#![allow(unused_imports)]

use wasmlib::*;
use wasmlib::host::*;

use crate::*;
use crate::keys::*;
use crate::structs::*;

#[derive(Clone, Copy)]
pub struct ArrayOfImmutableDonation {
	pub(crate) obj_id: i32,
}

impl ArrayOfImmutableDonation {
    pub fn length(&self) -> i32 {
        get_length(self.obj_id)
    }

	pub fn get_donation(&self, index: i32) -> ImmutableDonation {
		ImmutableDonation { obj_id: self.obj_id, key_id: Key32(index) }
	}
}

#[derive(Clone, Copy)]
pub struct ImmutableDonateWithFeedbackState {
    pub(crate) id: i32,
}

impl ImmutableDonateWithFeedbackState {
    pub fn log(&self) -> ArrayOfImmutableDonation {
		let arr_id = get_object_id(self.id, idx_map(IDX_STATE_LOG), TYPE_ARRAY | TYPE_BYTES);
		ArrayOfImmutableDonation { obj_id: arr_id }
	}

    pub fn max_donation(&self) -> ScImmutableInt64 {
		ScImmutableInt64::new(self.id, idx_map(IDX_STATE_MAX_DONATION))
	}

    pub fn total_donation(&self) -> ScImmutableInt64 {
		ScImmutableInt64::new(self.id, idx_map(IDX_STATE_TOTAL_DONATION))
	}
}

#[derive(Clone, Copy)]
pub struct ArrayOfMutableDonation {
	pub(crate) obj_id: i32,
}

impl ArrayOfMutableDonation {
    pub fn clear(&self) {
        clear(self.obj_id);
    }

    pub fn length(&self) -> i32 {
        get_length(self.obj_id)
    }

	pub fn get_donation(&self, index: i32) -> MutableDonation {
		MutableDonation { obj_id: self.obj_id, key_id: Key32(index) }
	}
}

#[derive(Clone, Copy)]
pub struct MutableDonateWithFeedbackState {
    pub(crate) id: i32,
}

impl MutableDonateWithFeedbackState {
    pub fn as_immutable(&self) -> ImmutableDonateWithFeedbackState {
		ImmutableDonateWithFeedbackState { id: self.id }
	}

    pub fn log(&self) -> ArrayOfMutableDonation {
		let arr_id = get_object_id(self.id, idx_map(IDX_STATE_LOG), TYPE_ARRAY | TYPE_BYTES);
		ArrayOfMutableDonation { obj_id: arr_id }
	}

    pub fn max_donation(&self) -> ScMutableInt64 {
		ScMutableInt64::new(self.id, idx_map(IDX_STATE_MAX_DONATION))
	}

    pub fn total_donation(&self) -> ScMutableInt64 {
		ScMutableInt64::new(self.id, idx_map(IDX_STATE_TOTAL_DONATION))
	}
}
