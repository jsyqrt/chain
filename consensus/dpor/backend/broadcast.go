package backend

import (
	"time"

	"bitbucket.org/cpchain/chain/commons/log"
	"bitbucket.org/cpchain/chain/types"
)

// BroadcastMinedBlock broadcasts generated block to committee
func (h *Handler) BroadcastMinedBlock(block *types.Block) {
	h.lock.Lock()
	defer h.lock.Unlock()

	committee := h.dialer.remoteValidators
	log.Debug("broadcast new generated block to commttee", "number", block.NumberU64())
	for addr, peer := range committee {
		log.Debug("broadcast new generated block to commttee", "addr", addr.Hex())
		peer.AsyncSendNewPendingBlock(block)
	}
}

// BroadcastPrepareSignedHeader broadcasts signed prepare header to remote committee
func (h *Handler) BroadcastPrepareSignedHeader(header *types.Header) {
	h.lock.Lock()
	defer h.lock.Unlock()

	committee := h.dialer.remoteValidators
	for _, peer := range committee {
		peer.AsyncSendPrepareSignedHeader(header)
	}
}

// BroadcastCommitSignedHeader broadcasts signed commit header to remote committee
func (h *Handler) BroadcastCommitSignedHeader(header *types.Header) {
	h.lock.Lock()
	defer h.lock.Unlock()

	committee := h.dialer.remoteValidators
	for _, peer := range committee {
		peer.AsyncSendCommitSignedHeader(header)
	}
}

// PendingBlockBroadcastLoop loops to broadcast blocks
func (h *Handler) PendingBlockBroadcastLoop() {
	futureTimer := time.NewTicker(time.Duration(h.dpor.ImpeachTimeout()))
	defer futureTimer.Stop()

	for {
		select {
		case pendingBlock := <-h.pendingBlockCh:

			log.Debug("proposed new pending block, broadcasting")

			ready := false

			for !ready {
				if h.Available() && len(h.dialer.remoteValidators) >= int(h.config.TermLen) {
					ready = true
				}
				time.Sleep(1 * time.Second)

				log.Debug("signer in dpor handler when broadcasting...")
				for addr := range h.dialer.remoteValidators {
					log.Debug("signer", "addr", addr.Hex())
				}
			}

			// broadcast mined pending block to remote signers
			go h.BroadcastMinedBlock(pendingBlock)

		case <-futureTimer.C:

			// check if still not received new block, if true, continue
			if h.ReadyToImpeach() {
				// get empty block
				if impeachBlock, err := h.dpor.CreateImpeachBlock(); err == nil {

					// broadcast the empty block
					// h.BroadcastGeneratedBlock(emptyBlock)
					_ = impeachBlock
					go h.BroadcastMinedBlock(impeachBlock) // TODO: @liuq BroadcastMinedBlock sends PREPREPARE message, which should be PREPARE message. We need to create a new function to let it calls.
				}
			}

		case <-h.quitSync:
			return
		}
	}
}
