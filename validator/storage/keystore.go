package storage

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts"
	"gitlab.waterfall.network/waterfall/protocol/gwat/accounts/keystore"
	"gitlab.waterfall.network/waterfall/protocol/gwat/common"
	"gitlab.waterfall.network/waterfall/protocol/gwat/console/prompt"
	"gitlab.waterfall.network/waterfall/protocol/gwat/log"
	"gitlab.waterfall.network/waterfall/protocol/gwat/node"
)

type VerifiersKeystore interface {
	UnlockAccounts(accounts [][]common.Address) error
	IsAddressUnlocked(address common.Address) bool
	UnlockAllAccounts()
}

type verifiersKeystore struct {
	keystore *keystore.KeyStore
	cfg      *node.VerifiersKeystoreConfig
}

func NewKeystore(cfg *node.VerifiersKeystoreConfig) VerifiersKeystore {
	keyStoreDir, err := cfg.KeyDir()
	if err != nil {
		return nil
	}

	return &verifiersKeystore{
		keystore: keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP),
		cfg:      cfg,
	}
}

// UnlockAccounts unlocks a specified accounts list.
func (k *verifiersKeystore) UnlockAccounts(accounts [][]common.Address) error {
	passwords, err := k.getPasswords()
	if err != nil {
		return err
	}
	keystoreAccounts := k.keystore.Accounts()

	for _, slotAccount := range accounts {
		for _, account := range slotAccount {
			// Find the position of the target account.
			position := findAccountPosition(keystoreAccounts, account.String())
			if position < 0 {
				continue
			}

			err = k.unlockAccount(account.String(), position, passwords)
			if err != nil {
				log.Error("unlock account failed", "error", err, "account", account.String())
			}
		}
	}

	return nil
}

func (k *verifiersKeystore) UnlockAllAccounts() {
	passwords, err := k.getPasswords()
	if err != nil {
		return
	}
	keystoreAccounts := k.keystore.Accounts()

	for _, account := range keystoreAccounts {
		// Find the position of the target account.
		position := findAccountPosition(keystoreAccounts, account.Address.String())
		if position < 0 {
			continue
		}

		err = k.unlockAccount(account.Address.String(), position, passwords)
		if err != nil {
			log.Error("unlock account failed", "error", err, "account", account.Address.String())
		}
	}
}

func (k *verifiersKeystore) IsAddressUnlocked(address common.Address) bool {
	return k.keystore.IsUnlocked(address)
}

// getPasswords returns a list of passwords from the password directory.
func (k *verifiersKeystore) getPasswords() ([]string, error) {
	dir, err := k.cfg.PasswordsFilePath()
	if err != nil {
		return nil, err
	}

	return makePasswordList(dir)
}

// MakeAddress converts an account specified directly as a hex encoded string or
// a key index in the key store to an internal account representation.
func (k *verifiersKeystore) makeAddress(account string) (accounts.Account, error) {
	// If the specified account is a valid address, return it
	if common.IsHexAddress(account) {
		return accounts.Account{Address: common.HexToAddress(account)}, nil
	}
	// Otherwise try to interpret the account as a keystore index
	index, err := strconv.Atoi(account)
	if err != nil || index < 0 {
		return accounts.Account{}, fmt.Errorf("invalid account address or index %q", account)
	}
	log.Warn("-------------------------------------------------------------------")
	log.Warn("Referring to accounts by order in the keystore folder is dangerous!")
	log.Warn("This functionality is deprecated and will be removed in the future!")
	log.Warn("Please use explicit addresses! (can search via `geth account list`)")
	log.Warn("-------------------------------------------------------------------")

	accs := k.keystore.Accounts()
	if len(accs) <= index {
		return accounts.Account{}, fmt.Errorf("index %d higher than number of accounts %d", index, len(accs))
	}
	return accs[index], nil
}

// GetPassPhraseWithList retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func getPassPhraseWithList(confirmation bool, index int, passwords []string) string {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if index < len(passwords) {
			return passwords[index]
		}
		return passwords[len(passwords)-1]
	}
	// Otherwise prompt the user for the password
	password := getPassPhrase(confirmation)
	return password
}

// GetPassPhrase displays the given text(prompt) to the user and requests some textual
// data to be entered, but one which must not be echoed out into the terminal.
// The method returns the input provided by the user.
func getPassPhrase(confirmation bool) string {
	password, err := prompt.Stdin.PromptPassword("Password: ")
	if err != nil {
		log.Error("Failed to read password:", "err", err)
	}
	if confirmation {
		confirm, err := prompt.Stdin.PromptPassword("Repeat password: ")
		if err != nil {
			log.Error("Failed to read password confirmation:", "err", err)
		}
		if password != confirm {
			log.Error("Passwords do not match")
		}
	}
	return password
}

// findAccountPosition finds the position of a target account in a list of accounts.
func findAccountPosition(accounts []accounts.Account, targetAddress string) int {
	for index, account := range accounts {
		if account.Address.String() == targetAddress {
			return index
		}
	}
	return -1 // account not found
}

// tries unlocking the specified account a few times.
func (k *verifiersKeystore) unlockAccount(address string, pos int, passwords []string) error {
	account, err := k.makeAddress(address)
	if err != nil {
		log.Error("Could not list accounts", "error", err)
		return err
	}

	if k.keystore.IsUnlocked(account.Address) {
		return nil
	}

	for trials := 0; trials < 3; trials++ {
		password := getPassPhraseWithList(false, pos, passwords)
		err = k.keystore.Unlock(account, password)
		if err != nil {
			log.Warn("Failed to unlock account, try again", "account", address, "error", err)
			continue
		}

		log.Info("Unlocked account", "address", account.Address.Hex())

		return nil
	}

	return nil
}

// MakePasswordList reads password lines from the file specified by the global --password flag.
func makePasswordList(path string) ([]string, error) {
	text, err := os.ReadFile(path)
	if err != nil {
		log.Error("Failed to read password file", "error", err)
		return nil, err
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines, err
}
