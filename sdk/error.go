package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/text/language"
)

// Existing CDS errors
// Note: the error id is useless except to ensure objects are different in map
var (
	ErrUnknownError                 = &Error{ID: 1, Status: http.StatusInternalServerError}
	ErrActionAlreadyUpdated         = &Error{ID: 2, Status: http.StatusBadRequest}
	ErrNoAction                     = &Error{ID: 3, Status: http.StatusNotFound}
	ErrActionLoop                   = &Error{ID: 4, Status: http.StatusBadRequest}
	ErrInvalidID                    = &Error{ID: 5, Status: http.StatusBadRequest}
	ErrInvalidProject               = &Error{ID: 6, Status: http.StatusBadRequest}
	ErrInvalidProjectKey            = &Error{ID: 7, Status: http.StatusBadRequest}
	ErrProjectHasPipeline           = &Error{ID: 8, Status: http.StatusConflict}
	ErrProjectHasApplication        = &Error{ID: 9, Status: http.StatusConflict}
	ErrUnauthorized                 = &Error{ID: 10, Status: http.StatusUnauthorized}
	ErrForbidden                    = &Error{ID: 11, Status: http.StatusForbidden}
	ErrPipelineNotFound             = &Error{ID: 12, Status: http.StatusBadRequest}
	ErrPipelineNotAttached          = &Error{ID: 13, Status: http.StatusBadRequest}
	ErrNoEnvironmentProvided        = &Error{ID: 14, Status: http.StatusBadRequest}
	ErrEnvironmentProvided          = &Error{ID: 15, Status: http.StatusBadRequest}
	ErrUnknownEnv                   = &Error{ID: 16, Status: http.StatusBadRequest}
	ErrEnvironmentExist             = &Error{ID: 17, Status: http.StatusConflict}
	ErrNoPipelineBuild              = &Error{ID: 18, Status: http.StatusNotFound}
	ErrApplicationNotFound          = &Error{ID: 19, Status: http.StatusNotFound}
	ErrGroupNotFound                = &Error{ID: 20, Status: http.StatusNotFound}
	ErrInvalidUsername              = &Error{ID: 21, Status: http.StatusBadRequest}
	ErrInvalidEmail                 = &Error{ID: 22, Status: http.StatusBadRequest}
	ErrGroupPresent                 = &Error{ID: 23, Status: http.StatusBadRequest}
	ErrInvalidName                  = &Error{ID: 24, Status: http.StatusBadRequest}
	ErrInvalidUser                  = &Error{ID: 25, Status: http.StatusBadRequest}
	ErrBuildArchived                = &Error{ID: 26, Status: http.StatusBadRequest}
	ErrNoEnvironment                = &Error{ID: 27, Status: http.StatusNotFound}
	ErrModelNameExist               = &Error{ID: 28, Status: http.StatusConflict}
	ErrNoWorkerModel                = &Error{ID: 29, Status: http.StatusNotFound}
	ErrNoProject                    = &Error{ID: 30, Status: http.StatusNotFound}
	ErrVariableExists               = &Error{ID: 31, Status: http.StatusConflict}
	ErrInvalidGroupPattern          = &Error{ID: 32, Status: http.StatusBadRequest}
	ErrGroupExists                  = &Error{ID: 33, Status: http.StatusConflict}
	ErrNotEnoughAdmin               = &Error{ID: 34, Status: http.StatusBadRequest}
	ErrInvalidProjectName           = &Error{ID: 35, Status: http.StatusBadRequest}
	ErrInvalidApplicationPattern    = &Error{ID: 36, Status: http.StatusBadRequest}
	ErrInvalidPipelinePattern       = &Error{ID: 37, Status: http.StatusBadRequest}
	ErrNotFound                     = &Error{ID: 38, Status: http.StatusNotFound}
	ErrNoWorkerModelCapa            = &Error{ID: 39, Status: http.StatusNotFound}
	ErrNoHook                       = &Error{ID: 40, Status: http.StatusNotFound}
	ErrNoAttachedPipeline           = &Error{ID: 41, Status: http.StatusNotFound}
	ErrNoReposManager               = &Error{ID: 42, Status: http.StatusNotFound}
	ErrNoReposManagerAuth           = &Error{ID: 43, Status: http.StatusUnauthorized}
	ErrNoReposManagerClientAuth     = &Error{ID: 44, Status: http.StatusForbidden}
	ErrRepoNotFound                 = &Error{ID: 45, Status: http.StatusNotFound}
	ErrSecretStoreUnreachable       = &Error{ID: 46, Status: http.StatusMethodNotAllowed}
	ErrSecretKeyFetchFailed         = &Error{ID: 47, Status: http.StatusMethodNotAllowed}
	ErrInvalidGoPath                = &Error{ID: 48, Status: http.StatusBadRequest}
	ErrCommitsFetchFailed           = &Error{ID: 49, Status: http.StatusNotFound}
	ErrInvalidSecretFormat          = &Error{ID: 50, Status: http.StatusInternalServerError}
	ErrUnknownTemplate              = &Error{ID: 51, Status: http.StatusNotFound}
	ErrNoPreviousSuccess            = &Error{ID: 52, Status: http.StatusNotFound}
	ErrNoEnvExecution               = &Error{ID: 53, Status: http.StatusForbidden}
	ErrSessionNotFound              = &Error{ID: 54, Status: http.StatusUnauthorized}
	ErrInvalidSecretValue           = &Error{ID: 55, Status: http.StatusBadRequest}
	ErrPipelineHasApplication       = &Error{ID: 56, Status: http.StatusBadRequest}
	ErrNoDirectSecretUse            = &Error{ID: 57, Status: http.StatusForbidden}
	ErrNoBranch                     = &Error{ID: 58, Status: http.StatusNotFound}
	ErrLDAPConn                     = &Error{ID: 59, Status: http.StatusInternalServerError}
	ErrServiceUnavailable           = &Error{ID: 60, Status: http.StatusServiceUnavailable}
	ErrParseUserNotification        = &Error{ID: 61, Status: http.StatusBadRequest}
	ErrNotSupportedUserNotification = &Error{ID: 62, Status: http.StatusBadRequest}
	ErrGroupNeedAdmin               = &Error{ID: 63, Status: http.StatusBadRequest}
	ErrGroupNeedWrite               = &Error{ID: 64, Status: http.StatusBadRequest}
	ErrNoVariable                   = &Error{ID: 65, Status: http.StatusNotFound}
	ErrPluginInvalid                = &Error{ID: 66, Status: http.StatusBadRequest}
	ErrConflict                     = &Error{ID: 67, Status: http.StatusConflict}
	ErrPipelineAlreadyAttached      = &Error{ID: 68, Status: http.StatusConflict}
	ErrApplicationExist             = &Error{ID: 69, Status: http.StatusConflict}
	ErrBranchNameNotProvided        = &Error{ID: 70, Status: http.StatusBadRequest}
	ErrInfiniteTriggerLoop          = &Error{ID: 71, Status: http.StatusBadRequest}
	ErrInvalidResetUser             = &Error{ID: 72, Status: http.StatusBadRequest}
	ErrUserConflict                 = &Error{ID: 73, Status: http.StatusBadRequest}
)

// SupportedLanguages on API errors
var SupportedLanguages = []language.Tag{
	language.AmericanEnglish,
	language.French,
}

var errorsAmericanEnglish = map[int]string{
	ErrUnknownError.ID:              "internal server error",
	ErrActionAlreadyUpdated.ID:      "action status already updated",
	ErrNoAction.ID:                  "action does not exist",
	ErrActionLoop.ID:                "action definition contains a recursive loop",
	ErrInvalidID.ID:                 "ID must be an integer",
	ErrInvalidProject.ID:            "project not provided",
	ErrInvalidProjectKey.ID:         "project key must contain only upper-case alphanumerical characters",
	ErrProjectHasPipeline.ID:        "project contains a pipeline",
	ErrProjectHasApplication.ID:     "project contains an application",
	ErrUnauthorized.ID:              "not authenticated",
	ErrForbidden.ID:                 "forbidden",
	ErrPipelineNotFound.ID:          "pipeline does not exist",
	ErrPipelineNotAttached.ID:       "pipeline is not attached to application",
	ErrNoEnvironmentProvided.ID:     "deployment and testing pipelines require an environnement",
	ErrEnvironmentProvided.ID:       "build pipeline are not compatible with environnment usage",
	ErrUnknownEnv.ID:                "unkown environment",
	ErrEnvironmentExist.ID:          "environment already exist",
	ErrNoPipelineBuild.ID:           "this pipeline build does not exist",
	ErrApplicationNotFound.ID:       "application does not exist",
	ErrGroupNotFound.ID:             "group does not exist",
	ErrInvalidUsername.ID:           "invalid username",
	ErrInvalidEmail.ID:              "invalid email",
	ErrGroupPresent.ID:              "group already present",
	ErrInvalidName.ID:               "invalid name",
	ErrInvalidUser.ID:               "invalid user or password",
	ErrBuildArchived.ID:             "Cannot restart this build because it has been archived",
	ErrNoEnvironment.ID:             "environment does not exist",
	ErrModelNameExist.ID:            "worker model name already used",
	ErrNoWorkerModel.ID:             "worker model does not exist",
	ErrNoProject.ID:                 "project does not exist",
	ErrVariableExists.ID:            "variable already exists",
	ErrInvalidGroupPattern.ID:       "group name must respect '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrGroupExists.ID:               "group already exists",
	ErrNotEnoughAdmin.ID:            "not enough group admin left",
	ErrInvalidProjectName.ID:        "project name must not be empty",
	ErrInvalidApplicationPattern.ID: "application name must respect '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrInvalidPipelinePattern.ID:    "pipeline name must respect '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrNotFound.ID:                  "page not found",
	ErrNoWorkerModelCapa.ID:         "capability not found",
	ErrNoHook.ID:                    "hook not found",
	ErrNoAttachedPipeline.ID:        "pipeline not attached to the application",
	ErrNoReposManager.ID:            "repositories manager not found",
	ErrNoReposManagerAuth.ID:        "CDS authentication error, please contact CDS administrator",
	ErrNoReposManagerClientAuth.ID:  "Repository manager authentication error, please check your credentials",
	ErrRepoNotFound.ID:              "repository not found",
	ErrSecretStoreUnreachable.ID:    "could not reach vault to fetch secret key",
	ErrSecretKeyFetchFailed.ID:      "error while fetching key from vault",
	ErrCommitsFetchFailed.ID:        "unable to retrieves commits",
	ErrInvalidGoPath.ID:             "invalid gopath",
	ErrInvalidSecretFormat.ID:       "cannot decrypt secret, invalid format",
	ErrUnknownTemplate.ID:           "unknown template",
	ErrSessionNotFound.ID:           "invalid session",
	ErrNoPreviousSuccess.ID:         "there is no previous success version for this pipeline",
	ErrNoEnvExecution.ID:            "you don't have execution right on this environment",
	ErrInvalidSecretValue.ID:        "secret value not specified",
	ErrPipelineHasApplication.ID:    "pipeline still used by an application",
	ErrNoDirectSecretUse.ID: `From now on, usage of 'password' parameter is not possible anymore.
Please use either Project, Application or Environment variable to store your secrets.
You can safely use them in a String or Text parameter`,
	ErrNoBranch.ID:                     "branch not found in repository",
	ErrLDAPConn.ID:                     "LDAP server connection error",
	ErrServiceUnavailable.ID:           "service currently unavailable or down for maintenance",
	ErrParseUserNotification.ID:        "unrecognized user notification settings",
	ErrNotSupportedUserNotification.ID: "unsupported user notification",
	ErrGroupNeedAdmin.ID:               "need at least 1 administrator",
	ErrGroupNeedWrite.ID:               "need at least 1 group with write permission",
	ErrNoVariable.ID:                   "variable not found",
	ErrPluginInvalid.ID:                "invalid plugin",
	ErrConflict.ID:                     "object conflict",
	ErrPipelineAlreadyAttached.ID:      "pipeline already attached to this application",
	ErrApplicationExist.ID:             "application already exist",
	ErrBranchNameNotProvided.ID:        "branchName parameter must be provided",
	ErrInfiniteTriggerLoop.ID:          "infinite trigger loop are forbidden",
	ErrInvalidResetUser.ID:             "invalid user or email",
	ErrUserConflict.ID:                 "this user already exist",
}

var errorsFrench = map[int]string{
	ErrUnknownError.ID:              "erreur interne",
	ErrActionAlreadyUpdated.ID:      "le status de l'action a déjà été mis à jour",
	ErrNoAction.ID:                  "l'action n'existe pas",
	ErrActionLoop.ID:                "la définition de l'action contient une boucle récursive",
	ErrInvalidID.ID:                 "l'ID doit être un nombre entier",
	ErrInvalidProject.ID:            "project manquant",
	ErrInvalidProjectKey.ID:         "la clef de project doit uniquement contenir des lettres majuscules et des chiffres",
	ErrProjectHasPipeline.ID:        "le project contient un pipeline",
	ErrProjectHasApplication.ID:     "le project contient une application",
	ErrUnauthorized.ID:              "authentification invalide",
	ErrForbidden.ID:                 "accès refusé",
	ErrPipelineNotFound.ID:          "le pipeline n'existe pas",
	ErrPipelineNotAttached.ID:       "le pipeline n'est pas lié à l'application",
	ErrNoEnvironmentProvided.ID:     "les pipelines de deploiement et de tests requiert un environnement",
	ErrEnvironmentProvided.ID:       "une pipeline de build ne necessite pas d'environnement",
	ErrUnknownEnv.ID:                "environnement inconnu",
	ErrEnvironmentExist.ID:          "l'environnement existe",
	ErrNoPipelineBuild.ID:           "ce build n'existe pas",
	ErrApplicationNotFound.ID:       "l'application n'existe pas",
	ErrGroupNotFound.ID:             "le groupe n'existe pas",
	ErrInvalidUsername.ID:           "nom d'utilisateur invalide",
	ErrInvalidEmail.ID:              "addresse email invalide",
	ErrGroupPresent.ID:              "le groupe est déjà présent",
	ErrInvalidName.ID:               "le nom est invalide",
	ErrInvalidUser.ID:               "mauvaise combinaison compte/mot de passe utilisateur",
	ErrBuildArchived.ID:             "impossible de relancer ce build car il a été archivé",
	ErrNoEnvironment.ID:             "l'environement n'existe pas",
	ErrModelNameExist.ID:            "le nom du modèle de worker est déjà utilisé",
	ErrNoWorkerModel.ID:             "le modèle de worker n'existe pas",
	ErrNoProject.ID:                 "le projet n'existe pas",
	ErrVariableExists.ID:            "la variable existe déjà",
	ErrInvalidGroupPattern.ID:       "nom de groupe invalide '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrGroupExists.ID:               "le groupe existe déjà",
	ErrNotEnoughAdmin.ID:            "pas assez d'admin restant",
	ErrInvalidProjectName.ID:        "nom de project vide non autorisé",
	ErrInvalidApplicationPattern.ID: "nom de l'application invalide '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrInvalidPipelinePattern.ID:    "nom du pipeline invalide '^[a-zA-Z0-9.-_-]{1,}$'",
	ErrNotFound.ID:                  "la page n'existe pas",
	ErrNoWorkerModelCapa.ID:         "la capacité n'existe pas",
	ErrNoHook.ID:                    "le hook n'existe pas",
	ErrNoAttachedPipeline.ID:        "le pipeline n'est pas lié à l'application",
	ErrNoReposManager.ID:            "le gestionnaire de dépôt n'existe pas",
	ErrNoReposManagerAuth.ID:        "connexion de CDS au gestionnaire de dépôt refusée, merci de contacter l'administrateur",
	ErrNoReposManagerClientAuth.ID:  "connexion au gestionnaire de dépôts refusée, vérifier vos accréditations",
	ErrRepoNotFound.ID:              "le dépôt n'existe pas",
	ErrSecretStoreUnreachable.ID:    "impossible de contacter vault",
	ErrSecretKeyFetchFailed.ID:      "erreur pendnat la récuperation de la clef de chiffrement",
	ErrCommitsFetchFailed.ID:        "impossible de retrouver les changements",
	ErrInvalidGoPath.ID:             "le gopath n'est pas valide",
	ErrInvalidSecretFormat.ID:       "impossibe de dechiffrer le secret, format invalide",
	ErrUnknownTemplate.ID:           "template inconnu",
	ErrSessionNotFound.ID:           "session invalide",
	ErrNoPreviousSuccess.ID:         "il n'y a aucune précédente version en succès pour ce pipeline",
	ErrNoEnvExecution.ID:            "vous n'avez pas les droits d'éxécution sur cet environnement",
	ErrInvalidSecretValue.ID:        "valeur du secret non specifiée",
	ErrPipelineHasApplication.ID:    "le pipeline est utilisé par une application",
	ErrNoDirectSecretUse.ID: `Désormais, l'utilisation du type de paramêtre 'password' n'est plus possible.
Merci d'utiliser les variable de Projet, Application ou Environnement pour stocker vos secrets.
Vous pouvez les utiliser sans problème dans un paramêtre de type String ou Text`,
	ErrNoBranch.ID:                     "la branche est introuvable dans le dépôt",
	ErrLDAPConn.ID:                     "erreur de connexion au serveur LDAP",
	ErrServiceUnavailable.ID:           "service temporairement indisponible ou en maintenance",
	ErrParseUserNotification.ID:        "notification non reconnue",
	ErrNotSupportedUserNotification.ID: "notification non supportée",
	ErrGroupNeedAdmin.ID:               "il faut au moins 1 administrateur",
	ErrGroupNeedWrite.ID:               "il faut au moins 1 groupe avec les droits d'écriture",
	ErrNoVariable.ID:                   "la variable n'existe pas",
	ErrPluginInvalid.ID:                "plugin non valide",
	ErrConflict.ID:                     "l'objet est en conflit",
	ErrPipelineAlreadyAttached.ID:      "le pipeline est déjà attaché à cette application",
	ErrApplicationExist.ID:             "une application du même nom existe déjà",
	ErrBranchNameNotProvided.ID:        "le paramètre branchName doit être envoyé",
	ErrInfiniteTriggerLoop.ID:          "création d'une boucle de trigger infinie interdite",
	ErrInvalidResetUser.ID:             "mauvaise combinaison compte/mail utilisateur",
	ErrUserConflict.ID:                 "cet utilisateur existe deja",
}

var matcher = language.NewMatcher(SupportedLanguages)

// NewError just set an error with a root cause
func NewError(target *Error, root error) *Error {
	target.Root = root
	return target
}

// ProcessError tries to recognize given error and return error message in a language matching Accepted-Language
func ProcessError(target error, al string) (string, int) {
	cdsErr, ok := target.(*Error)
	if !ok {
		return errorsAmericanEnglish[ErrUnknownError.ID], ErrUnknownError.Status
	}
	acceptedLanguages, _, err := language.ParseAcceptLanguage(al)
	if err != nil {
		return errorsAmericanEnglish[ErrUnknownError.ID], ErrUnknownError.Status
	}

	tag, _, _ := matcher.Match(acceptedLanguages...)
	var msg string
	switch tag {
	case language.French:
		msg, ok = errorsFrench[cdsErr.ID]
		break
	case language.AmericanEnglish:
		msg, ok = errorsAmericanEnglish[cdsErr.ID]
		break
	default:
		msg, ok = errorsAmericanEnglish[cdsErr.ID]
		break
	}

	if !ok {
		return errorsAmericanEnglish[ErrUnknownError.ID], ErrUnknownError.Status
	}
	if cdsErr.Root != nil {
		cdsErrRoot, ok := cdsErr.Root.(*Error)
		var rootMsg string
		if ok {
			rootMsg, _ = ProcessError(cdsErrRoot, al)
		} else {
			rootMsg = fmt.Sprintf("[%T] %s", cdsErr.Root, cdsErr.Root.Error())
		}

		msg = fmt.Sprintf("%s (caused by: %s)", msg, rootMsg)
	}
	return msg, cdsErr.Status
}

// Exit func display an error message on stderr and exit 1
func Exit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

// Error type
type Error struct {
	ID      int    `json:"-"`
	Status  int    `json:"-"`
	Message string `json:"message"`
	Root    error  `json:"-"`
}

// DecodeError return an Error struct from json
func DecodeError(data []byte) error {
	var e Error

	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil
	}

	if e.Message == "" {
		return nil
	}

	return e
}

func (e Error) String() string {
	if e.Message == "" {
		msg, ok := errorsAmericanEnglish[e.ID]
		if ok {
			return msg
		}
	}

	return e.Message
}

func (e Error) Error() string {
	return e.String()
}
