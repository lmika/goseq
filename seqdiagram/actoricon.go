package seqdiagram

import (
    "errors"

    "bitbucket.org/lmika/goseq/seqdiagram/graphbox"
)

// An actor icon
type ActorIcon interface {
    // Get the appropriate graphbox icon
    graphboxIcon() graphbox.Icon
}

// Error returned if the icon cannot be found
var EIconNotFound = errors.New("Icon not found")


// Lookup an actor icon based on it's name.  If the actor icon cannot be
// found, an EIconNotFound error is returned
func LookupActorIcon(name string) (ActorIcon, error) {
    // Lookup builtin icons
    if builtinIcon, hasBuiltinIcon := builtinIcons[name] ; hasBuiltinIcon {
        return builtinIcon, nil
    }

    return nil, EIconNotFound
}

// A build-in actor icon
type builtinActorIcon struct {
    icon           graphbox.Icon
}

func (bai *builtinActorIcon) graphboxIcon() graphbox.Icon {
    return bai.icon
}

// The set of built-in icons
var builtinIcons = map[string]ActorIcon {
    "human": &builtinActorIcon{graphbox.StickPersonIcon(1)},
    "cylinder": &builtinActorIcon{graphbox.CylinderIcon(1)},
}
