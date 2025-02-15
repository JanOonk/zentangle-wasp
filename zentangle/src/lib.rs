// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

#![allow(dead_code)]
#![allow(unused_imports)]

use zentangle::*;
use wasmlib::*;
use wasmlib::host::*;

use crate::consts::*;
use crate::events::*;
use crate::keys::*;
use crate::params::*;
use crate::results::*;
use crate::state::*;

mod consts;
mod contract;
mod events;
mod keys;
mod params;
mod results;
mod state;
mod structs;
mod zentangle;

#[no_mangle]
fn on_load() {
    let exports = ScExports::new();
    exports.add_func(FUNC_CREATE_GAME,         func_create_game_thunk);
    exports.add_func(FUNC_END_GAME,            func_end_game_thunk);
    exports.add_func(FUNC_INIT,                func_init_thunk);
    exports.add_func(FUNC_REQUEST_PLAY,        func_request_play_thunk);
    exports.add_func(FUNC_SEND_TAGS,           func_send_tags_thunk);
    exports.add_func(FUNC_SET_OWNER,           func_set_owner_thunk);
    exports.add_func(FUNC_WITHDRAW,            func_withdraw_thunk);
    exports.add_view(VIEW_GET_OWNER,           view_get_owner_thunk);
    exports.add_view(VIEW_GET_PLAYER_BETS,     view_get_player_bets_thunk);
    exports.add_view(VIEW_GET_PLAYER_INFO,     view_get_player_info_thunk);
    exports.add_view(VIEW_GET_PLAYS_PER_IMAGE, view_get_plays_per_image_thunk);
    exports.add_view(VIEW_GET_RESULTS,         view_get_results_thunk);

    unsafe {
        for i in 0..KEY_MAP_LEN {
            IDX_MAP[i] = get_key_id_from_string(KEY_MAP[i]);
        }
    }
}

pub struct CreateGameContext {
	events:  ZentangleEvents,
	params: ImmutableCreateGameParams,
	state: MutablezentangleState,
}

fn func_create_game_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcCreateGame");
	let f = CreateGameContext {
		events:  ZentangleEvents {},
		params: ImmutableCreateGameParams {
			id: OBJ_ID_PARAMS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.description().exists(), "missing mandatory description");
	ctx.require(f.params.number_of_images().exists(), "missing mandatory numberOfImages");
	func_create_game(ctx, &f);
	ctx.log("zentangle.funcCreateGame ok");
}

pub struct EndGameContext {
	events:  ZentangleEvents,
	params: ImmutableEndGameParams,
	state: MutablezentangleState,
}

fn func_end_game_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcEndGame");
	let f = EndGameContext {
		events:  ZentangleEvents {},
		params: ImmutableEndGameParams {
			id: OBJ_ID_PARAMS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	func_end_game(ctx, &f);
	ctx.log("zentangle.funcEndGame ok");
}

pub struct InitContext {
	events:  ZentangleEvents,
	params: ImmutableInitParams,
	state: MutablezentangleState,
}

fn func_init_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcInit");
	let f = InitContext {
		events:  ZentangleEvents {},
		params: ImmutableInitParams {
			id: OBJ_ID_PARAMS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	func_init(ctx, &f);
	ctx.log("zentangle.funcInit ok");
}

pub struct RequestPlayContext {
	events:  ZentangleEvents,
	results: MutableRequestPlayResults,
	state: MutablezentangleState,
}

fn func_request_play_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcRequestPlay");
	let f = RequestPlayContext {
		events:  ZentangleEvents {},
		results: MutableRequestPlayResults {
			id: OBJ_ID_RESULTS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	func_request_play(ctx, &f);
	ctx.log("zentangle.funcRequestPlay ok");
}

pub struct SendTagsContext {
	events:  ZentangleEvents,
	params: ImmutableSendTagsParams,
	state: MutablezentangleState,
}

fn func_send_tags_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcSendTags");
	let f = SendTagsContext {
		events:  ZentangleEvents {},
		params: ImmutableSendTagsParams {
			id: OBJ_ID_PARAMS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.input_json().exists(), "missing mandatory inputJson");
	func_send_tags(ctx, &f);
	ctx.log("zentangle.funcSendTags ok");
}

pub struct SetOwnerContext {
	events:  ZentangleEvents,
	params: ImmutableSetOwnerParams,
	state: MutablezentangleState,
}

fn func_set_owner_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcSetOwner");

	// current owner of this smart contract
	let access = ctx.state().get_agent_id("owner");
	ctx.require(access.exists(), "access not set: owner");
	ctx.require(ctx.caller() == access.value(), "no permission");

	let f = SetOwnerContext {
		events:  ZentangleEvents {},
		params: ImmutableSetOwnerParams {
			id: OBJ_ID_PARAMS,
		},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.owner().exists(), "missing mandatory owner");
	func_set_owner(ctx, &f);
	ctx.log("zentangle.funcSetOwner ok");
}

pub struct WithdrawContext {
	events:  ZentangleEvents,
	state: MutablezentangleState,
}

fn func_withdraw_thunk(ctx: &ScFuncContext) {
	ctx.log("zentangle.funcWithdraw");

	// current owner of this smart contract
	let access = ctx.state().get_agent_id("owner");
	ctx.require(access.exists(), "access not set: owner");
	ctx.require(ctx.caller() == access.value(), "no permission");

	let f = WithdrawContext {
		events:  ZentangleEvents {},
		state: MutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	func_withdraw(ctx, &f);
	ctx.log("zentangle.funcWithdraw ok");
}

pub struct GetOwnerContext {
	results: MutableGetOwnerResults,
	state: ImmutablezentangleState,
}

fn view_get_owner_thunk(ctx: &ScViewContext) {
	ctx.log("zentangle.viewGetOwner");
	let f = GetOwnerContext {
		results: MutableGetOwnerResults {
			id: OBJ_ID_RESULTS,
		},
		state: ImmutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	view_get_owner(ctx, &f);
	ctx.log("zentangle.viewGetOwner ok");
}

pub struct GetPlayerBetsContext {
	results: MutableGetPlayerBetsResults,
	state: ImmutablezentangleState,
}

fn view_get_player_bets_thunk(ctx: &ScViewContext) {
	ctx.log("zentangle.viewGetPlayerBets");
	let f = GetPlayerBetsContext {
		results: MutableGetPlayerBetsResults {
			id: OBJ_ID_RESULTS,
		},
		state: ImmutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	view_get_player_bets(ctx, &f);
	ctx.log("zentangle.viewGetPlayerBets ok");
}

pub struct GetPlayerInfoContext {
	params: ImmutableGetPlayerInfoParams,
	results: MutableGetPlayerInfoResults,
	state: ImmutablezentangleState,
}

fn view_get_player_info_thunk(ctx: &ScViewContext) {
	ctx.log("zentangle.viewGetPlayerInfo");
	let f = GetPlayerInfoContext {
		params: ImmutableGetPlayerInfoParams {
			id: OBJ_ID_PARAMS,
		},
		results: MutableGetPlayerInfoResults {
			id: OBJ_ID_RESULTS,
		},
		state: ImmutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.player_address().exists(), "missing mandatory playerAddress");
	view_get_player_info(ctx, &f);
	ctx.log("zentangle.viewGetPlayerInfo ok");
}

pub struct GetPlaysPerImageContext {
	params: ImmutableGetPlaysPerImageParams,
	results: MutableGetPlaysPerImageResults,
	state: ImmutablezentangleState,
}

fn view_get_plays_per_image_thunk(ctx: &ScViewContext) {
	ctx.log("zentangle.viewGetPlaysPerImage");
	let f = GetPlaysPerImageContext {
		params: ImmutableGetPlaysPerImageParams {
			id: OBJ_ID_PARAMS,
		},
		results: MutableGetPlaysPerImageResults {
			id: OBJ_ID_RESULTS,
		},
		state: ImmutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.image_id().exists(), "missing mandatory imageId");
	view_get_plays_per_image(ctx, &f);
	ctx.log("zentangle.viewGetPlaysPerImage ok");
}

pub struct GetResultsContext {
	params: ImmutableGetResultsParams,
	results: MutableGetResultsResults,
	state: ImmutablezentangleState,
}

fn view_get_results_thunk(ctx: &ScViewContext) {
	ctx.log("zentangle.viewGetResults");
	let f = GetResultsContext {
		params: ImmutableGetResultsParams {
			id: OBJ_ID_PARAMS,
		},
		results: MutableGetResultsResults {
			id: OBJ_ID_RESULTS,
		},
		state: ImmutablezentangleState {
			id: OBJ_ID_STATE,
		},
	};
	ctx.require(f.params.image_id().exists(), "missing mandatory imageId");
	view_get_results(ctx, &f);
	ctx.log("zentangle.viewGetResults ok");
}
